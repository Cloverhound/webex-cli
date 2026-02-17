package appconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Compiled-in defaults (set via ldflags or hardcoded).
var (
	DefaultClientID     = ""
	DefaultClientSecret = ""
	DefaultScopes       = "spark:all spark-admin:all"
)

type UserInfo struct {
	DisplayName string `json:"display_name"`
	OrgID       string `json:"org_id"`
	OrgName     string `json:"org_name"`
}

type Config struct {
	DefaultUser  string              `json:"default_user"`
	Users        map[string]UserInfo `json:"users"`
	ClientID     string              `json:"client_id,omitempty"`
	ClientSecret string              `json:"client_secret,omitempty"`
	Scopes       string              `json:"scopes,omitempty"`

	path string // file path, not serialized
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".webex-cli")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

// Load reads config from disk. Returns empty config if file doesn't exist.
func Load() (*Config, error) {
	cfg := &Config{
		Users: make(map[string]UserInfo),
		path:  configPath(),
	}

	data, err := os.ReadFile(cfg.path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Users == nil {
		cfg.Users = make(map[string]UserInfo)
	}
	return cfg, nil
}

// Save writes config to disk with secure permissions.
func (c *Config) Save() error {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(c.path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// EffectiveClientID returns the user-configured client ID or the compiled-in default.
func (c *Config) EffectiveClientID() string {
	if c.ClientID != "" {
		return c.ClientID
	}
	return DefaultClientID
}

// EffectiveClientSecret returns the user-configured client secret or the compiled-in default.
func (c *Config) EffectiveClientSecret() string {
	if c.ClientSecret != "" {
		return c.ClientSecret
	}
	return DefaultClientSecret
}

func (c *Config) EffectiveScopes() string {
	if c.Scopes != "" {
		return c.Scopes
	}
	return DefaultScopes
}

func (c *Config) SetDefaultUser(email string) {
	c.DefaultUser = email
}

func (c *Config) AddUser(email, displayName, orgID, orgName string) {
	c.Users[email] = UserInfo{
		DisplayName: displayName,
		OrgID:       orgID,
		OrgName:     orgName,
	}
}

func (c *Config) RemoveUser(email string) {
	delete(c.Users, email)
	if c.DefaultUser == email {
		c.DefaultUser = ""
	}
}

func (c *Config) UserEmails() []string {
	emails := make([]string, 0, len(c.Users))
	for e := range c.Users {
		emails = append(emails, e)
	}
	return emails
}
