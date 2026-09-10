package commands

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newDiscountCodesUpdateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "update"}
	cmd.Flags().String("code", "", "Discount code")
	cmd.Flags().String("kind", "", "Code kind (percentage, fixed)")
	cmd.Flags().Int("value", 0, "Discount value")
	cmd.Flags().String("meeting-ids", "", "Comma-separated meeting SIDs")
	cmd.Flags().String("expires-at", "", "Expiration date (YYYY-MM-DD)")
	cmd.Flags().String("json-file", "", "Path to JSON file with discount code data")
	return cmd
}

func TestDiscountCodeUpdateBody_OmitsFlagsTheCallerNeverSet(t *testing.T) {
	cmd := newDiscountCodesUpdateTestCmd()
	if err := cmd.Flags().Set("code", "SAVE20"); err != nil {
		t.Fatalf("set code: %v", err)
	}

	body, err := discountCodeUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	want := map[string]interface{}{"code": "SAVE20"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v; an unset --value must not send 0 and blank the discount", body, want)
	}
}

func TestDiscountCodeUpdateBody_SendsEveryFlagTheCallerSet(t *testing.T) {
	cmd := newDiscountCodesUpdateTestCmd()
	for flag, value := range map[string]string{
		"code":        "SAVE20",
		"kind":        "percentage",
		"value":       "20",
		"meeting-ids": "abc123x,def456y",
		"expires-at":  "2026-12-31",
	} {
		if err := cmd.Flags().Set(flag, value); err != nil {
			t.Fatalf("set %s: %v", flag, err)
		}
	}

	body, err := discountCodeUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	want := map[string]interface{}{
		"code":        "SAVE20",
		"kind":        "percentage",
		"value":       20,
		"meeting_ids": []string{"abc123x", "def456y"},
		"expires_at":  "2026-12-31",
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestDiscountCodeUpdateBody_KeepsAnExplicitZeroValue(t *testing.T) {
	cmd := newDiscountCodesUpdateTestCmd()
	if err := cmd.Flags().Set("value", "0"); err != nil {
		t.Fatalf("set value: %v", err)
	}

	body, err := discountCodeUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}
	if body["value"] != 0 {
		t.Errorf("body[value] = %v, want 0 to reach the API so it can reject it", body["value"])
	}
}

func TestDiscountCodeUpdateBody_LetsFlagsOverrideTheJSONFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "code.json")
	contents := `{"code":"FROMFILE","redeems_limit":5}`
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	cmd := newDiscountCodesUpdateTestCmd()
	if err := cmd.Flags().Set("json-file", path); err != nil {
		t.Fatalf("set json-file: %v", err)
	}
	if err := cmd.Flags().Set("code", "FROMFLAG"); err != nil {
		t.Fatalf("set code: %v", err)
	}

	body, err := discountCodeUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	if body["code"] != "FROMFLAG" {
		t.Errorf("body[code] = %v, want the flag to win over the file", body["code"])
	}
	if body["redeems_limit"] != float64(5) {
		t.Errorf("body[redeems_limit] = %v, want the file to supply fields that have no flag", body["redeems_limit"])
	}
}
