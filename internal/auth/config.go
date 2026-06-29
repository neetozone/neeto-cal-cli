package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	configDir = ".config/neetocal"
	authFile  = "auth.json"
)

type Credentials struct {
	Subdomain    string `json:"subdomain"`
	Email        string `json:"email"`
	SessionToken string `json:"session_token"`
}

// Store is the persisted collection of all logged-in subdomains.
type Store struct {
	Credentials []Credentials `json:"credentials"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Could not determine home directory: %w", err)
	}
	return filepath.Join(home, configDir), nil
}

func authFilePath() (string, error) {
	dir, err := configPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, authFile), nil
}

// LoadStore reads auth.json. Returns an empty Store (no error) when the file
// is absent. Transparently migrates the legacy single-object format.
func LoadStore() (*Store, error) {
	path, err := authFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Store{}, nil
		}
		return nil, fmt.Errorf("Could not read credentials: %w", err)
	}

	// New format: {"credentials": [...]}.
	var store Store
	if err := json.Unmarshal(data, &store); err == nil && store.Credentials != nil {
		return &store, nil
	}

	// Legacy format: a single Credentials object at the top level.
	var legacy Credentials
	if err := json.Unmarshal(data, &legacy); err != nil {
		return nil, fmt.Errorf("Invalid credentials file: %w", err)
	}
	if legacy.SessionToken == "" {
		return &Store{}, nil
	}
	return &Store{Credentials: []Credentials{legacy}}, nil
}

// SaveStore writes the store to disk, creating the config directory as needed.
// If the store holds no credentials the file is removed instead.
func SaveStore(store *Store) error {
	path, err := authFilePath()
	if err != nil {
		return err
	}

	if len(store.Credentials) == 0 {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("Could not remove credentials: %w", err)
		}
		return nil
	}

	dir, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("Could not create config directory: %w", err)
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("Could not serialize credentials: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("Could not write credentials: %w", err)
	}

	return nil
}

func (s *Store) Find(subdomain string) (*Credentials, bool) {
	for i := range s.Credentials {
		if s.Credentials[i].Subdomain == subdomain {
			return &s.Credentials[i], true
		}
	}
	return nil, false
}

// Upsert replaces the entry with the same subdomain, else appends a new one.
func (s *Store) Upsert(creds Credentials) {
	for i := range s.Credentials {
		if s.Credentials[i].Subdomain == creds.Subdomain {
			s.Credentials[i] = creds
			return
		}
	}
	s.Credentials = append(s.Credentials, creds)
}

// Remove deletes the entry with the given subdomain. Returns true if an entry
// was removed.
func (s *Store) Remove(subdomain string) bool {
	for i := range s.Credentials {
		if s.Credentials[i].Subdomain == subdomain {
			s.Credentials = append(s.Credentials[:i], s.Credentials[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) Subdomains() []string {
	out := make([]string, len(s.Credentials))
	for i, c := range s.Credentials {
		out[i] = c.Subdomain
	}
	return out
}

// SelectCredentials loads the store and returns the credentials a command
// should use, applying these rules:
//   - subdomain != "": return the matching entry, error if not present.
//   - subdomain == "" and store has one entry: return it (default).
//   - subdomain == "" and store is empty: "not authenticated" error.
//   - subdomain == "" and store has many entries: require --subdomain.
func SelectCredentials(subdomain string) (*Credentials, error) {
	store, err := LoadStore()
	if err != nil {
		return nil, err
	}
	if len(store.Credentials) == 0 {
		return nil, fmt.Errorf("Not authenticated. Run 'neetocal login' to authenticate.")
	}
	if subdomain != "" {
		creds, ok := store.Find(subdomain)
		if !ok {
			return nil, fmt.Errorf("Not authenticated for %q. Authenticated subdomains: %s.",
				subdomain, strings.Join(store.Subdomains(), ", "))
		}
		return creds, nil
	}
	if len(store.Credentials) == 1 {
		return &store.Credentials[0], nil
	}
	return nil, fmt.Errorf("Multiple subdomains authenticated (%s); specify --subdomain.",
		strings.Join(store.Subdomains(), ", "))
}
