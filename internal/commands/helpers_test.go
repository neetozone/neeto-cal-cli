package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newRequiredFlagsTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "create",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          func(cmd *cobra.Command, args []string) error { return nil },
	}
	cmd.Flags().String("name", "", "Name")
	cmd.Flags().String("time-zone", "", "Time zone")
	cmd.Flags().String("json-file", "", "Path to JSON file")
	markFlagsRequired(cmd, "name", "time-zone")
	allowJSONFileToSatisfyRequiredFlags(cmd)
	return cmd
}

func TestMarkFlagsRequired_FailsWithoutFlags(t *testing.T) {
	cmd := newRequiredFlagsTestCmd()
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() expected required flag error")
	}
	if !strings.Contains(err.Error(), "required flag(s)") {
		t.Errorf("error = %q, want required flag(s) error", err.Error())
	}
	if !strings.Contains(err.Error(), "name") || !strings.Contains(err.Error(), "time-zone") {
		t.Errorf("error = %q, want it to mention name and time-zone", err.Error())
	}
}

func TestMarkFlagsRequired_PassesWithFlags(t *testing.T) {
	cmd := newRequiredFlagsTestCmd()
	cmd.SetArgs([]string{"--name", "Work Hours", "--time-zone", "America/New_York"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
}

func TestMarkFlagsRequired_AppendsRequiredToUsage(t *testing.T) {
	cmd := newRequiredFlagsTestCmd()

	if usage := cmd.Flags().Lookup("name").Usage; !strings.Contains(usage, "(required)") {
		t.Errorf("name usage = %q, want it to contain (required)", usage)
	}
	if usage := cmd.Flags().Lookup("json-file").Usage; strings.Contains(usage, "(required)") {
		t.Errorf("json-file usage = %q, want it to not contain (required)", usage)
	}
}

func TestAllowJSONFileToSatisfyRequiredFlags_AllKeysInFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	if err := os.WriteFile(path, []byte(`{"name":"Work Hours","time_zone":"America/New_York"}`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cmd := newRequiredFlagsTestCmd()
	cmd.SetArgs([]string{"--json-file", path})

	if err := cmd.Execute(); err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
}

func TestAllowJSONFileToSatisfyRequiredFlags_EmptyValuesInFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	if err := os.WriteFile(path, []byte(`{"name":"","time_zone":null}`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cmd := newRequiredFlagsTestCmd()
	cmd.SetArgs([]string{"--json-file", path})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() expected required flag error for empty values")
	}
	if !strings.Contains(err.Error(), "name") || !strings.Contains(err.Error(), "time-zone") {
		t.Errorf("error = %q, want it to mention name and time-zone", err.Error())
	}
}

func TestAllowJSONFileToSatisfyRequiredFlags_ChainsExistingPreRunE(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	if err := os.WriteFile(path, []byte(`{"name":"Work Hours","time_zone":"America/New_York"}`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cmd := &cobra.Command{
		Use:           "create",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          func(cmd *cobra.Command, args []string) error { return nil },
	}
	cmd.Flags().String("name", "", "Name")
	cmd.Flags().String("time-zone", "", "Time zone")
	cmd.Flags().String("json-file", "", "Path to JSON file")
	markFlagsRequired(cmd, "name", "time-zone")

	existingRan := false
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		existingRan = true
		return nil
	}
	allowJSONFileToSatisfyRequiredFlags(cmd)

	cmd.SetArgs([]string{"--json-file", path})
	if err := cmd.Execute(); err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
	if !existingRan {
		t.Error("existing PreRunE was not chained")
	}
}

func TestAllowJSONFileToSatisfyRequiredFlags_MissingKeyInFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	if err := os.WriteFile(path, []byte(`{"name":"Work Hours"}`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cmd := newRequiredFlagsTestCmd()
	cmd.SetArgs([]string{"--json-file", path})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() expected required flag error")
	}
	if !strings.Contains(err.Error(), "time-zone") {
		t.Errorf("error = %q, want it to mention time-zone", err.Error())
	}
	if strings.Contains(err.Error(), `"name"`) {
		t.Errorf("error = %q, want it to not mention name", err.Error())
	}
}

func TestReadJSONFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	if err := os.WriteFile(path, []byte(`{"name":"test","count":42}`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, err := readJSONFile(path)
	if err != nil {
		t.Fatalf("readJSONFile() error = %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("name = %v, want test", result["name"])
	}
	if result["count"] != float64(42) {
		t.Errorf("count = %v, want 42", result["count"])
	}
}

func TestReadJSONFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{not valid`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := readJSONFile(path)
	if err == nil {
		t.Error("readJSONFile() expected error for invalid JSON")
	}
}

func TestReadJSONFile_NotFound(t *testing.T) {
	_, err := readJSONFile("/nonexistent/file.json")
	if err == nil {
		t.Error("readJSONFile() expected error for missing file")
	}
}
