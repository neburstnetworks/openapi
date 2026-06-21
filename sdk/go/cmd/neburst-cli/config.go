package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type Account struct {
	BaseURL string `json:"base_url"`
	KeyID   string `json:"key_id"`
	Secret  string `json:"secret"`
}

type Config struct {
	Active   string             `json:"active"`
	Accounts map[string]Account `json:"accounts"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".neburst", "config.json")
}

func loadConfig() *Config {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return &Config{Accounts: make(map[string]Account)}
	}
	var cfg Config
	if json.Unmarshal(data, &cfg) != nil {
		return &Config{Accounts: make(map[string]Account)}
	}
	if cfg.Accounts == nil {
		cfg.Accounts = make(map[string]Account)
	}
	return &cfg
}

func (c *Config) save() error {
	p := configPath()
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

func (c *Config) addAccount(name string, acct Account) {
	c.Accounts[name] = acct
}

func (c *Config) removeAccount(name string) {
	delete(c.Accounts, name)
	if c.Active == name {
		c.Active = ""
	}
}

func (c *Config) setActive(name string) bool {
	if _, ok := c.Accounts[name]; !ok {
		return false
	}
	c.Active = name
	return true
}

func (c *Config) activeAccount() (Account, bool) {
	acct, ok := c.Accounts[c.Active]
	return acct, ok
}

func (c *Config) sortedNames() []string {
	names := make([]string, 0, len(c.Accounts))
	for n := range c.Accounts {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func maskSecret(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

func parseCombinedKey(input string) (keyID, secret string, ok bool) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", "", false
	}
	var parsed struct {
		KeyID  string `json:"key_id"`
		Secret string `json:"secret"`
	}
	if json.Unmarshal(data, &parsed) == nil && parsed.KeyID != "" && parsed.Secret != "" {
		return parsed.KeyID, parsed.Secret, true
	}
	return "", "", false
}
