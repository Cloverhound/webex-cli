package config

import "testing"

const (
	testOrgUUID = "67b6b831-2fb5-448a-b5f5-b06e5e6253c2"
	// The ID Webex issues for testOrgUUID: base64 without padding.
	testOrgWebexID = "Y2lzY29zcGFyazovL3VzL09SR0FOSVpBVElPTi82N2I2YjgzMS0yZmI1LTQ0OGEtYjVmNS1iMDZlNWU2MjUzYzI"
)

func TestEncodeOrgIDMatchesWebexIssuedID(t *testing.T) {
	if got := EncodeOrgID(testOrgUUID); got != testOrgWebexID {
		t.Errorf("EncodeOrgID = %s, want %s", got, testOrgWebexID)
	}
}

func TestSetOrgIDNormalizesPaddedIDs(t *testing.T) {
	t.Cleanup(func() { SetOrgID("") })
	for _, in := range []string{testOrgUUID, testOrgWebexID, testOrgWebexID + "="} {
		SetOrgID(in)
		if OrgID() != testOrgUUID || OrgIDBase64() != testOrgWebexID {
			t.Errorf("SetOrgID(%s): uuid=%s base64=%s", in, OrgID(), OrgIDBase64())
		}
	}
}
