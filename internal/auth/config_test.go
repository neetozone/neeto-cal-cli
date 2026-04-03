package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadCredentials(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	creds := &Credentials{
		Subdomain:    "acme",
		Email:        "alice@example.com",
		SessionToken: "tok_abc123",
	}

	if err := SaveCredentials(creds); err != nil {
		t.Fatalf("SaveCredentials() error = %v", err)
	}

	loaded, err := LoadCredentials()
	if err != nil {
		t.Fatalf("LoadCredentials() error = %v", err)
	}

	if loaded.Subdomain != creds.Subdomain {
		t.Errorf("Subdomain = %q, want %q", loaded.Subdomain, creds.Subdomain)
	}
	if loaded.Email != creds.Email {
		t.Errorf("Email = %q, want %q", loaded.Email, creds.Email)
	}
	if loaded.SessionToken != creds.SessionToken {
		t.Errorf("SessionToken = %q, want %q", loaded.SessionToken, creds.SessionToken)
	}
}

func TestLoadCredentials_FileNotFound(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	_, err := LoadCredentials()
	if err == nil {
		t.Error("LoadCredentials() expected error when file does not exist")
	}
}

func TestLoadCredentials_EmptySessionToken(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	dir, _ := configPath()
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, authFile), []byte(`{"subdomain":"acme","email":"alice@example.com","session_token":""}`), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := LoadCredentials()
	if err == nil {
		t.Error("LoadCredentials() expected error for empty session token")
	}
}

func TestLoadCredentials_InvalidJSON(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	dir, _ := configPath()
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, authFile), []byte(`{not valid json`), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := LoadCredentials()
	if err == nil {
		t.Error("LoadCredentials() expected error for invalid JSON")
	}
}

func TestClearCredentials(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	creds := &Credentials{
		Subdomain:    "acme",
		Email:        "alice@example.com",
		SessionToken: "tok_abc123",
	}
	if err := SaveCredentials(creds); err != nil {
		t.Fatalf("SaveCredentials() error = %v", err)
	}

	if err := ClearCredentials(); err != nil {
		t.Fatalf("ClearCredentials() error = %v", err)
	}

	_, err := LoadCredentials()
	if err == nil {
		t.Error("LoadCredentials() expected error after ClearCredentials()")
	}
}

func TestClearCredentials_NoFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := ClearCredentials(); err != nil {
		t.Errorf("ClearCredentials() error = %v, want nil when file doesn't exist", err)
	}
}

func TestSaveCredentials_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	creds := &Credentials{
		Subdomain:    "acme",
		Email:        "alice@example.com",
		SessionToken: "tok_abc123",
	}

	if err := SaveCredentials(creds); err != nil {
		t.Fatalf("SaveCredentials() error = %v", err)
	}

	path := filepath.Join(tmpDir, configDir, authFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("SaveCredentials() did not create auth file")
	}
}

func TestSaveCredentials_FilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	creds := &Credentials{
		Subdomain:    "acme",
		Email:        "alice@example.com",
		SessionToken: "tok_abc123",
	}
	if err := SaveCredentials(creds); err != nil {
		t.Fatalf("SaveCredentials() error = %v", err)
	}

	path := filepath.Join(tmpDir, configDir, authFile)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("auth file permissions = %o, want 0600", perm)
	}
}
