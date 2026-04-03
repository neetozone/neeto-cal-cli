package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
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

func LoadCredentials() (*Credentials, error) {
	path, err := authFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not logged in. Run 'neetocal login' to authenticate")
		}
		return nil, fmt.Errorf("could not read credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("invalid credentials file: %w", err)
	}

	if creds.SessionToken == "" {
		return nil, fmt.Errorf("not logged in. Run 'neetocal login' to authenticate")
	}

	return &creds, nil
}

func SaveCredentials(creds *Credentials) error {
	dir, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("could not serialize credentials: %w", err)
	}

	path := filepath.Join(dir, authFile)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("could not write credentials: %w", err)
	}

	return nil
}

func ClearCredentials() error {
	path, err := authFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not remove credentials: %w", err)
	}

	return nil
}
