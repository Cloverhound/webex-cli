package webexapi

import "testing"

func TestDefaultHostsAreUS(t *testing.T) {
	h, err := NewHosts("", "")
	if err != nil {
		t.Fatal(err)
	}
	want := Hosts{
		Webex:            "https://webexapis.com/v1",
		CC:               "https://api.wxcc-us1.cisco.com",
		Analytics:        "https://analytics.webexapis.com",
		AnalyticsCalling: "https://analytics-calling.webexapis.com/v1",
		AIAgent:          "https://studio.aiagent-us1.cisco.com",
	}
	if h != want {
		t.Fatalf("got %+v", h)
	}
}

func TestContactCenterDataCenterSelectsCCAndStudioHosts(t *testing.T) {
	h, err := NewHosts("EU2", "eu")
	if err != nil {
		t.Fatal(err)
	}
	if h.CC != "https://api.wxcc-eu2.cisco.com" || h.AIAgent != "https://studio.aiagent-eu2.cisco.com" {
		t.Fatalf("got %+v", h)
	}
	if h.AnalyticsCalling != "https://analytics-calling-eu.webexapis.com/v1" {
		t.Fatalf("analytics calling = %s", h.AnalyticsCalling)
	}
}

func TestUnknownRegionsAreRefused(t *testing.T) {
	if _, err := NewHosts("mars1", ""); err == nil {
		t.Fatal("unknown data center accepted")
	}
	if _, err := NewHosts("", "mars"); err == nil {
		t.Fatal("unknown data region accepted")
	}
}

func TestOrgIDsConvertBothWays(t *testing.T) {
	uuid := "6ab1d2c4-0000-4000-8000-00000000c1d0"
	encoded := OrgIDFromUUID(uuid)
	if encoded != "Y2lzY29zcGFyazovL3VzL09SR0FOSVpBVElPTi82YWIxZDJjNC0wMDAwLTQwMDAtODAwMC0wMDAwMDAwMGMxZDA" {
		t.Fatalf("encoded = %s", encoded)
	}
	if OrgUUID(encoded) != uuid || OrgUUID(encoded+"=") != uuid || OrgUUID(uuid) != uuid {
		t.Fatal("decode failed")
	}
	if OrgIDFromUUID(encoded) != encoded || OrgUUID("other") != "other" || OrgIDFromUUID("") != "" {
		t.Fatal("non-matching input changed")
	}
}
