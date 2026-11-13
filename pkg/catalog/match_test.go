package catalog

import (
	"net/http"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// Every managed write is found from its own operation, and no two share one, so a
// recorded request names exactly one type.
func TestEveryManagedWriteIsFoundFromItsOperation(t *testing.T) {
	seen := map[string]string{}
	for _, typ := range Managed() {
		for action, r := range map[Action]resource.Route{Create: typ.Routes.Create, Update: typ.Routes.Update, Delete: typ.Routes.Delete} {
			if r.Op == nil {
				continue
			}
			key := string(r.Op.Host) + " " + r.Op.Method + " " + r.Op.Path
			if other, dup := seen[key]; dup {
				t.Errorf("%s and %s share %s", other, typ.Name, key)
			}
			seen[key] = typ.Name
			w, ok := WriteFor(r.Op.Host, r.Op.Method, r.Op.Path)
			if !ok || w.Type.Name != typ.Name || w.Action != action {
				t.Errorf("%s %s: found %v %s", typ.Name, action, w.Type.Name, w.Action)
			}
			got, _, ok := ResolveTarget(typ.Area, r.Template(), action == Create)
			if !ok || got.Name != typ.Name {
				t.Errorf("%s %s: target resolves to %q", typ.Name, action, got.Name)
			}
		}
	}
}

func TestWriteForRefusesOtherOperations(t *testing.T) {
	if _, ok := WriteFor(webexops.HostWebex, http.MethodPost, "/telephony/calls/answer"); ok {
		t.Error("call control is not a catalog write")
	}
	hg, _ := ByName("hunt_group")
	if _, ok := WriteFor(webexops.HostWebex, http.MethodPatch, hg.Routes.Update.Op.Path); ok {
		t.Error("the method must match too")
	}
	aa, _ := ByName("auto_attendant")
	if !aa.NoCreate || !aa.NoDelete {
		t.Fatal("the auto attendant is update-only")
	}
}

// A Contact Center queue is updated and deleted at different routes; both name it.
func TestAQueuesDeleteRouteResolvesForAnUpdate(t *testing.T) {
	q, _ := ByName("cc_queue")
	if q.Routes.Delete.Template() == q.Routes.Update.Template() {
		t.Skip("the queue's routes no longer differ")
	}
	got, r, ok := ResolveTarget(resource.ContactCenter, "/contact-service-queue/{queueId}", false)
	if !ok || got.Name != "cc_queue" {
		t.Fatalf("resolved %q", got.Name)
	}
	if key := KeyOf(r, "/contact-service-queue/{queueId}", map[string]string{"queueId": "q1"}); key.ID != "q1" {
		t.Errorf("key = %+v", key)
	}
}

func TestKeyOfPairsPlaceholdersByPosition(t *testing.T) {
	hg, _ := ByName("hunt_group")
	route := "/telephony/config/locations/{loc}/huntGroups/{hg}"
	key := KeyOf(hg.Routes.Update, route, map[string]string{"loc": "L", "hg": "H"})
	if key != (resource.Key{Parent: "L", ID: "H"}) {
		t.Errorf("key = %+v", key)
	}
}
