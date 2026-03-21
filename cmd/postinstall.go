package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const skillURL = "https://raw.githubusercontent.com/Cloverhound/webex-cli/main/skill/SKILL.md"

const coworkName = "Claude Cowork"

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
	RunE: func(cmd *cobra.Command, args []string) error {
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
				if err := setupUnixPath(installDir); err != nil {
					return err
				}
			}
			fmt.Println()
		}

		// Step 2: Agent skill installation
		if err := setupAgentSkills(); err != nil {
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

func setupUnixPath(installDir string) error {
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

	var choice string
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

func setupAgentSkills() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("detecting home directory: %w", err)
	}

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

	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Install Webex skill for AI coding agents?").
				Description(
					"The Webex skill lets AI coding agents (Claude Code, Codex, Cursor)\n"+
						"query and manage your Webex environment using natural language.\n"+
						"Detected agents are pre-selected.",
				).
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

	fmt.Print("Downloading Webex skill...")
	skillContent, err := downloadSkill()
	if err != nil {
		fmt.Println(" failed")
		return fmt.Errorf("downloading skill: %w", err)
	}
	fmt.Println(" ok")

	for _, name := range selected {
		for _, p := range agentPlatforms {
			if p.Name != name {
				continue
			}
			if p.Name == coworkName {
				zipPath := filepath.Join(home, "Downloads", "webex-cli-skill.zip")
				if err := buildSkillZip(zipPath, "webex-cli", skillContent); err != nil {
					fmt.Printf("  %s: failed (%v)\n", p.Name, err)
				} else {
					fmt.Printf("  %s: saved to %s\n", p.Name, zipPath)
					printCoworkInstructions(zipPath)
				}
			} else {
				dest := filepath.Join(home, p.SkillDir, "SKILL.md")
				if err := installSkill(dest, skillContent); err != nil {
					fmt.Printf("  %s: failed (%v)\n", p.Name, err)
				} else {
					fmt.Printf("  %s: installed to %s\n", p.Name, dest)
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

func downloadSkill() ([]byte, error) {
	resp, err := http.Get(skillURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
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

func buildSkillZip(dest, skillName string, skillContent []byte) error {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	f, err := w.Create(skillName + "/SKILL.md")
	if err != nil {
		return fmt.Errorf("creating zip entry: %w", err)
	}
	if _, err := f.Write(skillContent); err != nil {
		return fmt.Errorf("writing zip entry: %w", err)
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
	rootCmd.AddCommand(postInstallCmd)
}
