package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

func checkSkillUpdates() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	type skillAction struct {
		Platform agentPlatform
		Path     string
		Status   string // "outdated", "new"
	}

	fmt.Print("Checking for skill updates...")
	latest, err := downloadSkill()
	if err != nil {
		fmt.Println(" failed")
		return err
	}
	fmt.Println(" ok")

	var actions []skillAction
	for _, p := range agentPlatforms {
		if p.SkillDir == "" {
			continue // skip Cowork (no filesystem path)
		}
		path := filepath.Join(home, p.SkillDir, "SKILL.md")
		current, err := os.ReadFile(path)
		if err == nil {
			if !bytes.Equal(current, latest) {
				actions = append(actions, skillAction{p, path, "outdated"})
			}
		} else if agentDetected(home, p) {
			actions = append(actions, skillAction{p, path, "new"})
		}
	}

	if len(actions) == 0 {
		fmt.Println("All installed skills are up to date.")
		return nil
	}

	// Step 1: Ask if user wants to review skill updates
	var proceed bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Agent skill updates available").
				Description(fmt.Sprintf("%d agent skill(s) can be updated or installed.", len(actions))).
				Affirmative("Review changes").
				Negative("Skip").
				Value(&proceed),
		),
	)

	if err := confirmForm.Run(); err != nil || !proceed {
		return nil
	}

	// Step 2: Show multiselect with specific skills
	var options []huh.Option[string]
	for _, a := range actions {
		label := a.Platform.Name
		if a.Status == "outdated" {
			label += "  (update available)"
		} else {
			label += "  (not yet installed)"
		}
		options = append(options, huh.NewOption(label, a.Platform.Name).Selected(true))
	}

	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Skill updates available").
				Description("The following agents can be updated or have the Webex skill\ninstalled. Deselect any you want to skip.").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return nil
	}

	if len(selected) == 0 {
		return nil
	}

	for _, name := range selected {
		for _, a := range actions {
			if a.Platform.Name == name {
				if err := installSkill(a.Path, latest); err != nil {
					fmt.Printf("  %s: failed (%v)\n", a.Platform.Name, err)
				} else if a.Status == "outdated" {
					fmt.Printf("  %s: updated %s\n", a.Platform.Name, a.Path)
				} else {
					fmt.Printf("  %s: installed to %s\n", a.Platform.Name, a.Path)
				}
			}
		}
	}

	return nil
}
