package commands

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newPrefillTokenTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create"}
	cmd.Flags().String("data", "", "")
	cmd.Flags().String("json-file", "", "")
	return cmd
}

func TestPrefillTokenBody_ParsesDataFlag(t *testing.T) {
	cmd := newPrefillTokenTestCmd()
	setFlags(t, cmd, map[string]string{"data": `{"name":"Oliver Smith"}`})

	body, err := prefillTokenBody(cmd)
	if err != nil {
		t.Fatalf("body: %v", err)
	}

	want := map[string]interface{}{"data": map[string]any{"name": "Oliver Smith"}}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestPrefillTokenBody_ReadsJSONFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "answers.json")
	if err := os.WriteFile(path, []byte(`{"email":"oliver@example.com"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := newPrefillTokenTestCmd()
	setFlags(t, cmd, map[string]string{"json-file": path})

	body, err := prefillTokenBody(cmd)
	if err != nil {
		t.Fatalf("body: %v", err)
	}

	want := map[string]interface{}{"data": map[string]any{"email": "oliver@example.com"}}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestPrefillTokenBody_RequiresDataOrJSONFile(t *testing.T) {
	_, err := prefillTokenBody(newPrefillTokenTestCmd())

	if err == nil || err.Error() != "pass --data or --json-file" {
		t.Errorf("err = %v, want pass --data or --json-file", err)
	}
}

func TestPrefillTokenBody_RejectsInvalidData(t *testing.T) {
	cmd := newPrefillTokenTestCmd()
	setFlags(t, cmd, map[string]string{"data": "nope"})

	if _, err := prefillTokenBody(cmd); err == nil {
		t.Error("err = nil, want invalid JSON error")
	}
}
