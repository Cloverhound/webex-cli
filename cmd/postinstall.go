package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Cloverhound/webex-cli/skill"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const coworkName = "Claude Cowork"

// setupOptions controls how post-install and skill updates answer questions.
type setupOptions struct {
	prompt    bool // ask with interactive forms
	assumeYes bool // without prompts, take the recommended action
	skills    bool // install or update agent skills
}

// newSetupOptions prompts only when a person can answer. These forms are drawn
// on stdout, so it must be a terminal too.
func newSetupOptions(assumeYes, noSkills bool) setupOptions {
	prompt := !assumeYes && canPrompt() && isTerminal(os.Stdout)
	return setupOptions{prompt: prompt, assumeYes: assumeYes, skills: !noSkills}
}

type agentPlatform struct {
	Name     string
	SkillDir string // relative to home dir
}

var agentPlatforms = []agentPlatform{
	{"Claude Code", filepath.Join(".claude", "skills", "webex-cli")},
	{coworkName, ""}, // Cowork requires ZIP upload via web UI
	{"OpenAI Codex", filepath.Join(".codex", "skills", "webex-cli")},
	{"Cursor", filepath.Join(".cursor", "skills", "webex-cli")},
}

var postInstallCmd = &cobra.Command{
	Use:   "post-install",
	Short: "Run post-installation setup (PATH, agent skills)",
	Long: `Adds the install directory to PATH and installs the Webex agent skill.

Without a terminal, or with --yes or $WEBEX_CLI_NONINTERACTIVE=1, nothing is
asked: the skill is installed for detected agents (Claude Code, Codex, Cursor).
Without --yes, the PATH line to add is printed instead of editing shell files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		assumeYes, _ := cmd.Flags().GetBool("yes")
		noSkills, _ := cmd.Flags().GetBool("no-skills")
		skillsOnly, _ := cmd.Flags().GetBool("skills-only")
		opts := newSetupOptions(assumeYes, noSkills)

		if skillsOnly {
			return checkSkillUpdates(opts)
		}

		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("detecting executable path: %w", err)
		}
		installDir := filepath.Dir(execPath)

		// Step 1: PATH setup
		if !dirInPath(installDir) {
			if runtime.GOOS == "windows" {
				if err := setupWindowsPath(installDir); err != nil {
					return err
				}
			} else {
				if err := setupUnixPath(installDir, opts); err != nil {
					return err
				}
			}
			fmt.Println()
		}

		// Step 2: Agent skill installation
		if err := setupAgentSkills(opts); err != nil {
			return err
		}

		return nil
	},
}

func dirInPath(dir string) bool {
	for _, p := range filepath.SplitList(os.Getenv("PATH")) {
		if p == dir {
			return true
		}
	}
	return false
}

func setupUnixPath(installDir string, opts setupOptions) error {
	shell := filepath.Base(os.Getenv("SHELL"))
	var rcFile string
	switch shell {
	case "zsh":
		rcFile = filepath.Join(os.Getenv("HOME"), ".zshrc")
	case "bash":
		rcFile = filepath.Join(os.Getenv("HOME"), ".bashrc")
	default:
		rcFile = filepath.Join(os.Getenv("HOME"), ".profile")
	}

	exportLine := fmt.Sprintf(`export PATH="%s:$PATH"`, installDir)

	if data, err := os.ReadFile(rcFile); err == nil {
		if strings.Contains(string(data), installDir) {
			fmt.Printf("%s already references %s — restart your terminal or run: source %s\n", rcFile, installDir, rcFile)
			return nil
		}
	}

	choice := "no"
	if opts.assumeYes {
		choice = "yes"
	} else if opts.prompt {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(fmt.Sprintf("%s is not in your PATH", installDir)).
					Description(
						fmt.Sprintf(
							"The webex binary was installed to %s, but your shell\n"+
								"can't find it yet. Adding it to %s will make the\n"+
								"\"webex\" command available in all new terminal sessions.",
							installDir, filepath.Base(rcFile),
						),
					).
					Options(
						huh.NewOption(fmt.Sprintf("Yes — add to %s", filepath.Base(rcFile)), "yes"),
						huh.NewOption("No — I'll do it myself", "no"),
					).
					Value(&choice),
			),
		)

		if err := form.Run(); err != nil {
			return nil
		}
	} else {
		fmt.Printf("%s is not in your PATH.\n", installDir)
	}

	if choice == "yes" {
		f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("opening %s: %w", rcFile, err)
		}
		defer f.Close()

		if data, err := os.ReadFile(rcFile); err == nil && len(data) > 0 && data[len(data)-1] != '\n' {
			f.WriteString("\n")
		}
		f.WriteString(exportLine + "\n")

		fmt.Printf("Added to %s. Restart your terminal or run: source %s\n", rcFile, rcFile)
	} else {
		fmt.Println()
		fmt.Println("To add it manually, run:")
		fmt.Println()
		fmt.Printf("  echo '%s' >> %s && source %s\n", exportLine, rcFile, rcFile)
		fmt.Println()
	}

	return nil
}

func setupWindowsPath(installDir string) error {
	fmt.Printf("Add %s to your system PATH to use the \"webex\" command.\n", installDir)
	return nil
}

func setupAgentSkills(opts setupOptions) error {
	if !opts.skills {
		return nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("detecting home directory: %w", err)
	}

	var selected []string
	if opts.prompt {
		var options []huh.Option[string]
		for _, p := range agentPlatforms {
			label := p.Name
			if p.Name == coworkName {
				label += "  (saves ZIP to ~/Downloads for manual upload)"
			}
			opt := huh.NewOption(label, p.Name)
			if agentDetected(home, p) {
				opt = opt.Selected(true)
			}
			options = append(options, opt)
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Install Webex skill for AI coding agents?").
					Description(
						"The Webex skill lets AI coding agents (Claude Code, Codex, Cursor)\n" +
							"query and manage your Webex environment using natural language.\n" +
							"Detected agents are pre-selected.",
					).
					Options(options...).
					Value(&selected),
			),
		)

		if err := form.Run(); err != nil {
			return nil
		}
	} else {
		selected = detectedSkillAgents(home)
	}

	if len(selected) == 0 {
		return nil
	}

	skillFiles, err := skill.Files()
	if err != nil {
		return fmt.Errorf("reading embedded skill: %w", err)
	}

	for _, name := range selected {
		for _, p := range agentPlatforms {
			if p.Name != name {
				continue
			}
			if p.Name == coworkName {
				zipPath := filepath.Join(home, "Downloads", "webex-cli-skill.zip")
				if err := buildSkillZip(zipPath, "webex-cli", skillFiles); err != nil {
					fmt.Printf("  %s: failed (%v)\n", p.Name, err)
				} else {
					fmt.Printf("  %s: saved to %s\n", p.Name, zipPath)
					printCoworkInstructions(zipPath)
				}
			} else {
				skillDir := filepath.Join(home, p.SkillDir)
				failed := false
				for subPath, content := range skillFiles {
					dest := filepath.Join(skillDir, filepath.FromSlash(subPath))
					if err := installSkill(dest, content); err != nil {
						fmt.Printf("  %s: failed installing %s (%v)\n", p.Name, subPath, err)
						failed = true
						break
					}
				}
				if !failed {
					fmt.Printf("  %s: installed to %s\n", p.Name, skillDir)
				}
			}
		}
	}

	return nil
}

func agentDetected(home string, p agentPlatform) bool {
	if p.Name == coworkName {
		switch runtime.GOOS {
		case "darwin":
			_, err := os.Stat("/Applications/Claude.app")
			return err == nil
		case "windows":
			_, err := os.Stat(filepath.Join(os.Getenv("LOCALAPPDATA"), "AnthropicClaude"))
			return err == nil
		}
		return false
	}
	topDir := filepath.Join(home, strings.SplitN(p.SkillDir, string(filepath.Separator), 2)[0])
	info, err := os.Stat(topDir)
	return err == nil && info.IsDir()
}

// detectedSkillAgents lists the detected agents that install from a skill
// directory. Cowork is left out because it needs a manual ZIP upload.
func detectedSkillAgents(home string) []string {
	var names []string
	for _, p := range agentPlatforms {
		if p.SkillDir != "" && agentDetected(home, p) {
			names = append(names, p.Name)
		}
	}
	return names
}

func printCoworkInstructions(zipPath string) {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("212")).
		Padding(1, 2).
		MarginLeft(2)

	heading := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	step := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	sub := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Italic(true)

	content := heading.Render("⚠ Action required: Upload skill to Claude Cowork") + "\n" +
		sub.Render("Unlike Claude Code, Cowork skills must be manually uploaded to the Claude Desktop app.") + "\n\n" +
		step.Render("1. Open Claude Desktop and switch to the Cowork tab") + "\n" +
		step.Render("2. Click Customize (left sidebar) → Skills → + → Upload a skill") + "\n" +
		step.Render("3. Select: "+zipPath)

	fmt.Println(box.Render(content))
}

func buildSkillZip(dest, skillName string, files map[string][]byte) error {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	for subPath, content := range files {
		f, err := w.Create(skillName + "/" + subPath)
		if err != nil {
			return fmt.Errorf("creating zip entry %s: %w", subPath, err)
		}
		if _, err := f.Write(content); err != nil {
			return fmt.Errorf("writing zip entry %s: %w", subPath, err)
		}
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("closing zip: %w", err)
	}

	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	return os.WriteFile(dest, buf.Bytes(), 0644)
}

func installSkill(dest string, content []byte) error {
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	return os.WriteFile(dest, content, 0644)
}

func init() {
	postInstallCmd.Flags().BoolP("yes", "y", false, "Do not prompt; add the PATH line and install skills for detected agents")
	postInstallCmd.Flags().Bool("no-skills", false, "Skip agent skill installation")
	postInstallCmd.Flags().Bool("skills-only", false, "Only check installed agent skills against this binary's skill")
	postInstallCmd.Flags().MarkHidden("skills-only")
	rootCmd.AddCommand(postInstallCmd)
}
