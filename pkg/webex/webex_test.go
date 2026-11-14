package webex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Cloverhound/webex-cli/pkg/catalog"
	"github.com/Cloverhound/webex-cli/pkg/fakeapi"
	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

func newClient(t *testing.T) (*Client, *fakeapi.Server) {
	t.Helper()
	fake := fakeapi.New()
	srv := httptest.NewServer(fake)
	t.Cleanup(srv.Close)
	c := &Client{
		API:     &webexapi.Client{Token: func(context.Context) (string, error) { return "tok", nil }},
		Hosts:   fakeapi.Hosts(srv.URL),
		OrgUUID: fakeapi.OrgUUID,
	}
	return c, fake
}

func typeOf(t *testing.T, name string) resource.Type {
	t.Helper()
	typ, err := catalog.ByName(name)
	if err != nil {
		t.Fatal(err)
	}
	return typ
}

func TestListReadsDetailsForPartialLists(t *testing.T) {
	c, _ := newClient(t)
	objs, err := c.List(context.Background(), typeOf(t, "hunt_group"))
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) != 3 {
		t.Fatalf("objects = %d", len(objs))
	}
	o := objs[0]
	if o.Key != (resource.Key{Parent: "loc1", ID: "hg1"}) || o.Label != "g-hg1" {
		t.Fatalf("object = %+v", o)
	}
	if o.Values["no_answer_rings"] != 5.0 || o.Values["location_id"] != "loc1" {
		t.Fatalf("values = %v", o.Values)
	}
	if _, ok := o.Raw["agents"]; !ok {
		t.Fatal("raw must keep unmodelled fields such as agents for display")
	}
}

func TestListUsesFullItemsFromCCListsAndReadsReferenceTypes(t *testing.T) {
	c, _ := newClient(t)
	teams, err := c.List(context.Background(), typeOf(t, "cc_team"))
	if err != nil || len(teams) != 2 || teams[0].Values["site_id"] != "site1" {
		t.Fatalf("teams=%+v err=%v", teams, err)
	}
	locations, err := c.List(context.Background(), typeOf(t, "location"))
	if err != nil || len(locations) != 1 || locations[0].Label != "Boston" {
		t.Fatalf("locations=%+v err=%v", locations, err)
	}
	users, err := c.List(context.Background(), typeOf(t, "cc_user"))
	if err != nil || len(users) != 1 || users[0].Label != "ada@example.com" {
		t.Fatalf("users=%+v err=%v", users, err)
	}
}

func TestListCoversEveryPage(t *testing.T) {
	c, fake := newClient(t)
	fake.Seed(250)
	objs, err := c.List(context.Background(), typeOf(t, "hunt_group"))
	if err != nil || len(objs) != 253 {
		t.Fatalf("objects=%d err=%v", len(objs), err)
	}
}

func TestDetailReadsRunConcurrentlyAndKeepOrder(t *testing.T) {
	c, fake := newClient(t)
	fake.Seed(97)
	fake.SetLatency(20 * time.Millisecond)
	start := time.Now()
	objs, err := c.List(context.Background(), typeOf(t, "hunt_group"))
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("took %s; one at a time would be about 2s", elapsed)
	}
	if len(objs) != 100 || objs[0].Key.ID != "hg1" || objs[99].Key.ID != "hg1096" {
		t.Fatalf("len=%d first=%s last=%s", len(objs), objs[0].Key.ID, objs[len(objs)-1].Key.ID)
	}
}

func TestListFailsOnASingleFailedDetailRead(t *testing.T) {
	c, fake := newClient(t)
	fake.AddFault(fakeapi.Fault{Method: "GET", Path: "huntGroups/hg2$", Mode: fakeapi.Fail500NoCommit})
	if _, err := c.List(context.Background(), typeOf(t, "hunt_group")); err == nil {
		t.Fatal("a partial read must not pass as a whole org")
	}
}

func TestReadMissingObjectIsNotFound(t *testing.T) {
	c, _ := newClient(t)
	_, err := c.Read(context.Background(), typeOf(t, "cc_team"), resource.Key{ID: "nope"})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}

func TestCreateReturnsKeyWithParent(t *testing.T) {
	c, fake := newClient(t)
	typ := typeOf(t, "hunt_group")
	key, err := c.Create(context.Background(), typ, resource.Values{"location_id": "loc1", "name": "new", "extension": "1777"})
	if err != nil {
		t.Fatal(err)
	}
	if key.Parent != "loc1" || key.ID == "" {
		t.Fatalf("key = %+v", key)
	}
	if fake.HuntGroup(key.ID)["name"] != "new" {
		t.Fatal("object was not created")
	}
}

func TestAValueWebexWouldRefuseIsNeverSent(t *testing.T) {
	c, fake := newClient(t)
	typ := typeOf(t, "cc_global_variable")
	_, err := c.Create(context.Background(), typ, resource.Values{"name": "zz-studio-live-var", "variable_type": "String"})
	if err == nil || !strings.Contains(err.Error(), "Global Variable name") {
		t.Fatalf("create err = %v", err)
	}
	err = c.Update(context.Background(), typ, resource.Key{ID: "gv1"}, map[string]any{"id": "gv1"}, resource.Values{"name": "zz bad"})
	if err == nil || !strings.Contains(err.Error(), "Use only letters A–Z, digits and _") {
		t.Fatalf("update err = %v", err)
	}
	if m := fake.Mutations(); len(m) != 0 {
		t.Fatalf("sent %v", m)
	}
}

func TestAnUpdateToAFixedFieldIsNeverSent(t *testing.T) {
	c, fake := newClient(t)
	typ := typeOf(t, "cc_skill")
	got, err := c.Read(context.Background(), typ, resource.Key{ID: "skill1"})
	if err != nil {
		t.Fatal(err)
	}
	err = c.Update(context.Background(), typ, got.Key, got.Raw, resource.Values{"name": "zz renamed"})
	if err == nil || err.Error() != "Skill name: cannot be changed after the object is created" {
		t.Fatalf("err = %v", err)
	}
	if Classify(err) != Rejected {
		t.Fatal("a write refused before sending changed nothing")
	}
	if m := fake.Mutations(); len(m) != 0 {
		t.Fatalf("sent %v", m)
	}
}

func TestDeltaWriteSendsOnlyChangedFields(t *testing.T) {
	c, fake := newClient(t)
	typ := typeOf(t, "hunt_group")
	key := resource.Key{Parent: "loc1", ID: "hg1"}
	if err := c.Update(context.Background(), typ, key, nil, resource.Values{"no_answer_rings": 9}); err != nil {
		t.Fatal(err)
	}
	h := fake.HuntGroup("hg1")
	policies := h["callPolicies"].(map[string]any)["noAnswer"].(map[string]any)
	if policies["numberOfRings"] != 9.0 || policies["systemMaxNumberOfRings"] != 20.0 {
		t.Fatalf("policies = %v", policies)
	}
	if h["directoryPin"] != "s3cret-unmodelled" {
		t.Fatal("unmodelled hidden field was lost")
	}
}

func TestFullWriteKeepsMembershipAndUnmodelledFields(t *testing.T) {
	c, fake := newClient(t)
	typ := typeOf(t, "cc_team")
	key := resource.Key{ID: "team1"}
	got, err := c.Read(context.Background(), typ, key)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Update(context.Background(), typ, key, got.Raw, resource.Values{"description": "edited"}); err != nil {
		t.Fatal(err)
	}
	team := fake.Team("team1")
	if team["description"] != "edited" || team["futureField"] != "keep-me" {
		t.Fatalf("team = %v", team)
	}
	if !reflect.DeepEqual(team["userIds"], []any{"user1"}) {
		t.Fatalf("membership changed: %v", team["userIds"])
	}
}

func TestDelete(t *testing.T) {
	c, fake := newClient(t)
	if err := c.Delete(context.Background(), typeOf(t, "cc_team"), resource.Key{ID: "team2"}); err != nil {
		t.Fatal(err)
	}
	if fake.Team("team2") != nil {
		t.Fatal("team still exists")
	}
}

func TestPerLocationListsCarryTheirLocation(t *testing.T) {
	c, _ := newClient(t)
	objs, err := c.List(context.Background(), typeOf(t, "schedule"))
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) != 1 || objs[0].Key.Parent != "loc1" {
		t.Fatalf("objs = %+v", objs)
	}
	if _, leaked := objs[0].Raw[perLocationKey]; leaked {
		t.Fatal("the internal location key must not stay in the raw object")
	}
}

func TestSingletonsBecomeOneObjectEach(t *testing.T) {
	c, _ := newClient(t)
	org, err := c.List(context.Background(), typeOf(t, "voicemail_settings"))
	if err != nil || len(org) != 1 || org[0].Key.ID != "voicemail_settings" || org[0].Label == "" {
		t.Fatalf("org=%+v err=%v", org, err)
	}
	perLocation, err := c.List(context.Background(), typeOf(t, "location_voicemail"))
	if err != nil || len(perLocation) != 1 || perLocation[0].Key.Parent != "loc1" {
		t.Fatalf("perLocation=%+v err=%v", perLocation, err)
	}
}

func TestAnIDFieldOtherThanIDIsHonoured(t *testing.T) {
	c, _ := newClient(t)
	objs, err := c.List(context.Background(), typeOf(t, "phone_number"))
	if err != nil || len(objs) != 1 || objs[0].Key.ID == "" {
		t.Fatalf("objs=%+v err=%v", objs, err)
	}
}

func TestAnExtensionOnlyNumberUsesItsExtensionAsIDAndLabel(t *testing.T) {
	c, _ := newClient(t)
	typ := typeOf(t, "phone_number")
	obj, err := c.fromListItem(context.Background(), typ, map[string]any{"extension": "1234"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if obj.Key.ID != "1234" || obj.Label != "1234" {
		t.Fatalf("object = %+v", obj)
	}
}

func TestPermissionErrorsAreUnavailableButServerErrorsAreNot(t *testing.T) {
	if !IsUnavailable(&webexapi.APIError{Status: 403}) || !IsUnavailable(&webexapi.APIError{Status: 404}) {
		t.Fatal("a 403 or 404 means the org cannot read this route")
	}
	if !IsUnavailable(&webexapi.APIError{Status: 412}) {
		t.Fatal("a 412 means the feature is not available for this org")
	}
	if IsUnavailable(&webexapi.APIError{Status: 500}) || IsUnavailable(errors.New("reset")) || IsUnavailable(nil) {
		t.Fatal("anything else is a failed read")
	}
}

func TestEveryTypeListsAndReadsBack(t *testing.T) {
	c, _ := newClient(t)
	for _, typ := range catalog.All() {
		objs, err := c.List(context.Background(), typ)
		if err != nil {
			t.Fatalf("%s: %v", typ.Name, err)
		}
		if len(objs) == 0 {
			t.Fatalf("%s: no objects", typ.Name)
		}
		if objs[0].Label == "" {
			t.Errorf("%s has no label", typ.Name)
		}
	}
}

func TestContactCenterFullWritesKeepUnmodelledStructure(t *testing.T) {
	c, fake := newClient(t)
	for _, tc := range []struct{ typ, id, collection, kept string }{
		{"cc_queue", "queue1", "cc_queue", "callDistributionGroups"},
		{"cc_business_hours", "hours1", "cc_business_hours", "workingHours"},
		{"cc_entry_point", "ep1", "cc_entry_point", "controlFlowScriptUrl"},
		{"cc_skill", "skill1", "cc_skill", "dynamicSkill"},
	} {
		typ := typeOf(t, tc.typ)
		key := resource.Key{ID: tc.id}
		got, err := c.Read(context.Background(), typ, key)
		if err != nil {
			t.Fatalf("%s: %v", tc.typ, err)
		}
		if err := c.Update(context.Background(), typ, key, got.Raw, resource.Values{"description": "edited"}); err != nil {
			t.Fatalf("%s: %v", tc.typ, err)
		}
		stored := fake.Object(tc.collection, tc.id)
		if stored["description"] != "edited" {
			t.Fatalf("%s = %v", tc.typ, stored)
		}
		if _, ok := stored[tc.kept]; !ok {
			t.Fatalf("%s lost %s: %v", tc.typ, tc.kept, stored)
		}
	}
}

func TestFlowsListOutsideTheOrganizationPath(t *testing.T) {
	c, fake := newClient(t)
	flows, err := c.List(context.Background(), typeOf(t, "cc_flow"))
	if err != nil || len(flows) != 1 {
		t.Fatalf("flows = %v, %v", flows, err)
	}
	want := "/" + fakeapi.OrgUUID + "/project/5e5c9ad6d61f870d6d778c1b/flows"
	if log := fake.Log(); len(log) == 0 || log[len(log)-1].Path != want {
		t.Fatalf("read %v, want %s", log, want)
	}
}

func TestADialNumberWithOnlyAnExtensionIsLabelledByIt(t *testing.T) {
	c, fake := newClient(t)
	fake.Put("cc_dial_number", map[string]any{"id": "dn2", "extension": "5000", "entryPointId": "ep1"})
	objs, err := c.List(context.Background(), typeOf(t, "cc_dial_number"))
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range objs {
		if o.Key.ID == "dn2" && o.Label != "5000" {
			t.Fatalf("label = %q", o.Label)
		}
	}
}

func TestEntryPointFlowAndVersionWrite(t *testing.T) {
	c, fake := newClient(t)
	typ := typeOf(t, "cc_entry_point")
	key := resource.Key{ID: "ep1"}
	got, err := c.Read(context.Background(), typ, key)
	if err != nil {
		t.Fatal(err)
	}
	if got.Values["flow_id"] != "cc_flow1" || got.Values["flow_version"] != "Live" {
		t.Fatalf("values = %v", got.Values)
	}
	if err := c.Update(context.Background(), typ, key, got.Raw, resource.Values{"flow_id": "flow2", "flow_version": "Latest"}); err != nil {
		t.Fatal(err)
	}
	stored := fake.Object("cc_entry_point", "ep1")
	if stored["flowId"] != "flow2" || stored["flowTagId"] != "Latest" {
		t.Fatalf("stored = %v", stored)
	}
}

func TestCallingDeltaKeepsNestedSettings(t *testing.T) {
	c, fake := newClient(t)
	typ := typeOf(t, "call_queue")
	key := resource.Key{Parent: "loc1", ID: "cq1"}
	if err := c.Update(context.Background(), typ, key, nil, resource.Values{"queue_size": 7, "policy": "REGULAR"}); err != nil {
		t.Fatal(err)
	}
	q := fake.Object("call_queue", "cq1")
	settings := q["queueSettings"].(map[string]any)
	policies := q["callPolicies"].(map[string]any)
	if settings["queueSize"] != 7.0 || settings["overflow"] == nil {
		t.Fatalf("settings = %v", settings)
	}
	if policies["policy"] != "REGULAR" || policies["callBounce"] == nil {
		t.Fatalf("policies = %v", policies)
	}
	if q["agents"] == nil {
		t.Fatal("agents were lost")
	}
}

func TestCCQueueDeleteUsesItsOwnRoute(t *testing.T) {
	c, fake := newClient(t)
	if err := c.Delete(context.Background(), typeOf(t, "cc_queue"), resource.Key{ID: "queue1"}); err != nil {
		t.Fatal(err)
	}
	if fake.Object("cc_queue", "queue1") != nil {
		t.Fatal("queue still exists")
	}
}

func TestOrgIDIsSentToCallingWhenSet(t *testing.T) {
	c, fake := newClient(t)
	c.OrgID = "base64org"
	if _, err := c.Read(context.Background(), typeOf(t, "hunt_group"), resource.Key{Parent: "loc1", ID: "hg1"}); err != nil {
		t.Fatal(err)
	}
	if len(fake.Log()) != 1 {
		t.Fatalf("log = %v", fake.Log())
	}
}

func TestANewLayoutCopiesTheGlobalLayoutsFile(t *testing.T) {
	c, fake := newClient(t)
	layouts := typeOf(t, "cc_desktop_layout")
	for _, o := range fake.Objects("cc_desktop_layout") {
		o["global"] = false
		fake.Put("cc_desktop_layout", o)
	}
	fake.Put("cc_desktop_layout", map[string]any{"id": "global1", "name": "Global", "global": true, "jsonFileContent": `{"layout":1}`})
	key, err := c.Create(context.Background(), layouts, resource.Values{"name": "New"})
	if err != nil {
		t.Fatal(err)
	}
	got := fake.Object("cc_desktop_layout", key.ID)
	if got["jsonFileContent"] != `{"layout":1}` || got["name"] != "New" || got["global"] != false {
		t.Fatalf("created = %v", got)
	}
}

// A create that copies from another object checks just before its POST, after the
// reads, so a stop that arrives during them still prevents the write.
func TestBeforeWriteRunsAfterTheCopyReadsAndCanStopTheCreate(t *testing.T) {
	c, fake := newClient(t)
	fake.Put("cc_desktop_layout", map[string]any{"id": "global1", "name": "Global", "global": true, "jsonFileContent": `{"layout":1}`})
	var readsBefore int
	stop := errors.New("stopped")
	ctx := WithBeforeWrite(context.Background(), func() error {
		readsBefore = len(fake.Log())
		return stop
	})
	if _, err := c.Create(ctx, typeOf(t, "cc_desktop_layout"), resource.Values{"name": "New"}); !errors.Is(err, stop) {
		t.Fatalf("err = %v", err)
	}
	if readsBefore == 0 || len(fake.Mutations()) != 0 {
		t.Fatalf("check ran after %d requests; writes %v", readsBefore, fake.Mutations())
	}
}

// The operation table records lists that take no paging parameters, which answer
// in one response however many items they hold.
func TestAListWithoutPagingParametersIsReadInOneRequest(t *testing.T) {
	var queries []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.RawQuery)
		if len(queries) > 3 { // stops a pager that pages a list that does not page
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var list []map[string]any
		for i := 0; i < 150; i++ {
			list = append(list, map[string]any{"id": fmt.Sprintf("p%d", i), "name": fmt.Sprintf("Playlist %d", i)})
		}
		json.NewEncoder(w).Encode(map[string]any{"playlists": list})
	}))
	defer srv.Close()
	c := &Client{API: &webexapi.Client{Token: func(context.Context) (string, error) { return "tok", nil }},
		Hosts: fakeapi.Hosts(srv.URL), OrgUUID: fakeapi.OrgUUID}

	objects, err := c.List(context.Background(), typeOf(t, "announcement_playlist"))
	if err != nil || len(objects) != 150 {
		t.Fatalf("got %d, err %v", len(objects), err)
	}
	if len(queries) != 1 || strings.Contains(queries[0], "start=") {
		t.Fatalf("queries = %v", queries)
	}
}

func TestClassify(t *testing.T) {
	if Classify(&webexapi.APIError{Status: 400}) != Rejected {
		t.Fatal("a 400 is a clear refusal")
	}
	if Classify(&webexapi.APIError{Status: 500}) != Unknown {
		t.Fatal("a 500 may follow a committed write")
	}
	if Classify(errors.New("connection reset")) != Unknown {
		t.Fatal("a transport failure may follow a committed write")
	}
}
