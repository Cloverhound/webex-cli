package localconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const dirName = ".webex-cli"
const fileName = "config.json"

type Config struct {
	User string `json:"user"`

	path string
}

// Load reads the local .webex-cli/config.json from the given directory.
// Returns nil (no error) if the file doesn't exist.
func Load(dir string) (*Config, error) {
	p := filepath.Join(dir, dirName, fileName)
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading local config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing local config: %w", err)
	}
	cfg.path = p
	return &cfg, nil
}

// Save writes the local .webex-cli/config.json in the given directory.
func Save(dir string, username string) error {
	configDir := filepath.Join(dir, dirName)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("creating local config dir: %w", err)
	}

	cfg := Config{User: username}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling local config: %w", err)
	}
	data = append(data, '\n')

	p := filepath.Join(configDir, fileName)
	if err := os.WriteFile(p, data, 0600); err != nil {
		return fmt.Errorf("writing local config: %w", err)
	}

	ensureGitignore(dir)
	return nil
}

func ensureGitignore(dir string) {
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return
	}

	gitignorePath := filepath.Join(dir, ".gitignore")
	entry := ".webex-cli/"

	data, err := os.ReadFile(gitignorePath)
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.TrimSpace(line) == entry {
				return
			}
		}
	}

	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	if len(data) > 0 && data[len(data)-1] != '\n' {
		f.WriteString("\n")
	}
	f.WriteString(entry + "\n")
}
