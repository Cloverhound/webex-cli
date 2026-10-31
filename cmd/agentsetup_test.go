package cmd

import (
	"bytes"
	"testing"

	"github.com/Cloverhound/webex-cli/skill"
)

func TestAgentSetupPrintsSetupSkill(t *testing.T) {
	var out bytes.Buffer
	agentSetupCmd.SetOut(&out)
	defer agentSetupCmd.SetOut(nil)

	if err := agentSetupCmd.RunE(agentSetupCmd, nil); err != nil {
		t.Fatal(err)
	}
	if out.String() != skill.Setup {
		t.Errorf("printed %d bytes, want the %d-byte setup skill", out.Len(), len(skill.Setup))
	}
}

func TestAgentSetupNeedsNoLogin(t *testing.T) {
	c, _, err := rootCmd.Find([]string{"agent-setup"})
	if err != nil {
		t.Fatal(err)
	}
	if c != agentSetupCmd || !skipAuth(c) {
		t.Errorf("agent-setup resolved to %q, skipAuth = %v", c.Name(), skipAuth(c))
	}
}
