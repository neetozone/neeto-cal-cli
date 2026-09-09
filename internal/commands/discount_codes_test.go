package commands

import (
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
	return cmd
}

func TestDiscountCodeUpdateBody_OmitsFlagsTheCallerNeverSet(t *testing.T) {
	cmd := newDiscountCodesUpdateTestCmd()
	if err := cmd.Flags().Set("code", "SAVE20"); err != nil {
		t.Fatalf("set code: %v", err)
	}

	body := discountCodeUpdateBody(cmd)

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

	body := discountCodeUpdateBody(cmd)

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

	if body := discountCodeUpdateBody(cmd); body["value"] != 0 {
		t.Errorf("body[value] = %v, want 0 to reach the API so it can reject it", body["value"])
	}
}
