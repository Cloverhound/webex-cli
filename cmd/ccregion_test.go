package cmd

import (
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/config"
)

const (
	euOrg = "11111111-1111-1111-1111-111111111111"
	usOrg = "22222222-2222-2222-2222-222222222222"
)

// isolate gives a test its own config file and resets what it changes.
func isolate(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("WEBEX_CC_REGION", "")
	t.Setenv("WEBEX_USER", "")
	t.Cleanup(func() {
		rootCmd.PersistentFlags().Set("organization", "")
		rootCmd.PersistentFlags().Set("cc-region", "")
		rootCmd.SetArgs(nil)
		config.SetOrgID("")
		if err := setCCRegion(""); err != nil {
			t.Error(err)
		}
	})
}

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	rootCmd.SetArgs(args)
	var err error
	out := captureStdout(t, func() { err = rootCmd.Execute() })
	return out, err
}

func TestConfigSetCCRegionSavesItForTheOrganization(t *testing.T) {
	isolate(t)
	if _, err := run(t, "config", "set", "cc-region", "eu1", "--organization", config.EncodeOrgID(euOrg)); err != nil {
		t.Fatal(err)
	}
	cfg, err := appconfig.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := cfg.CCRegions[euOrg]; got != "eu1" {
		t.Errorf("saved region for %s = %q, want eu1 (saved: %v)", euOrg, got, cfg.CCRegions)
	}
	if got, ok := cfg.CCRegions[usOrg]; ok {
		t.Errorf("another organization got region %q", got)
	}
}

func TestConfigSetCCRegionUsesTheDefaultUsersOrganization(t *testing.T) {
	isolate(t)
	cfg, _ := appconfig.Load()
	cfg.DefaultUser = "me@example.com"
	cfg.Users["me@example.com"] = appconfig.UserInfo{OrgID: euOrg}
	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}
	if _, err := run(t, "config", "set", "cc-region", "anz1"); err != nil {
		t.Fatal(err)
	}
	cfg, _ = appconfig.Load()
	if got := cfg.CCRegions[euOrg]; got != "anz1" {
		t.Errorf("saved region = %q, want anz1", got)
	}
}

func TestConfigSetCCRegionNeedsAnOrganization(t *testing.T) {
	isolate(t)
	if _, err := run(t, "config", "set", "cc-region", "eu1"); err == nil || !strings.Contains(err.Error(), "organization") {
		t.Errorf("err = %v, want one asking for an organization", err)
	}
}

func TestConfigSetCCRegionRefusesAnUnknownRegion(t *testing.T) {
	isolate(t)
	if _, err := run(t, "config", "set", "cc-region", "eu9", "--organization", euOrg); err == nil || !strings.Contains(err.Error(), "eu9") {
		t.Errorf("err = %v, want one naming eu9", err)
	}
}

func TestConfigGetCCRegionShowsTheOrganizationsRegion(t *testing.T) {
	isolate(t)
	cfg, _ := appconfig.Load()
	cfg.CCRegions = map[string]string{euOrg: "eu2"}
	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}
	out, err := run(t, "config", "get", "cc-region", "--organization", euOrg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "eu2") || !strings.Contains(out, euOrg) {
		t.Errorf("output = %q, want eu2 for %s", out, euOrg)
	}
	out, err = run(t, "config", "get", "cc-region", "--organization", usOrg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "us1 (default)") {
		t.Errorf("output = %q, want us1 (default)", out)
	}
}

func TestTheSavedRegionSelectsTheOrganizationsHosts(t *testing.T) {
	isolate(t)
	cfg := &appconfig.Config{CCRegions: map[string]string{euOrg: "eu1"}}

	config.SetOrgID(euOrg)
	if err := resolveCCRegion("", cfg); err != nil {
		t.Fatal(err)
	}
	if config.CcBaseURL != "https://api.wxcc-eu1.cisco.com" || config.AiAgentBaseURL != "https://studio.aiagent-eu1.cisco.com" {
		t.Errorf("hosts = %s, %s; want eu1's", config.CcBaseURL, config.AiAgentBaseURL)
	}

	config.SetOrgID(usOrg)
	if err := resolveCCRegion("", cfg); err != nil {
		t.Fatal(err)
	}
	if config.CcBaseURL != "https://api.wxcc-us1.cisco.com" || config.AiAgentBaseURL != "https://studio.aiagent-us1.cisco.com" {
		t.Errorf("hosts = %s, %s; want us1's for an organization with no saved region", config.CcBaseURL, config.AiAgentBaseURL)
	}
}

func TestTheFlagAndEnvironmentOverrideTheSavedRegion(t *testing.T) {
	isolate(t)
	cfg := &appconfig.Config{CCRegions: map[string]string{euOrg: "eu1"}}
	config.SetOrgID(euOrg)

	t.Setenv("WEBEX_CC_REGION", "ca1")
	if err := resolveCCRegion("", cfg); err != nil {
		t.Fatal(err)
	}
	if config.CcBaseURL != "https://api.wxcc-ca1.cisco.com" {
		t.Errorf("with WEBEX_CC_REGION=ca1, CC host = %s", config.CcBaseURL)
	}

	if err := resolveCCRegion("JP1", cfg); err != nil {
		t.Fatal(err)
	}
	if config.CcBaseURL != "https://api.wxcc-jp1.cisco.com" {
		t.Errorf("with --cc-region JP1, CC host = %s", config.CcBaseURL)
	}

	if err := resolveCCRegion("mars", cfg); err == nil {
		t.Error("an unknown region was accepted")
	}
}
