package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestPaymentFilterParams_SendsFilters(t *testing.T) {
	cmd := &cobra.Command{Use: "list"}
	addPaymentFilterFlags(cmd, "payments", "")
	setFlags(t, cmd, map[string]string{
		"from-time":   "2026-10-01",
		"to-time":     "2026-11-01",
		"status":      "declined",
		"type":        "package_purchase",
		"meeting-sid": "abc123x",
	})

	params := paymentFilterParams(cmd)

	for param, want := range map[string]string{
		"from_time":   "2026-10-01",
		"to_time":     "2026-11-01",
		"status":      "declined",
		"type":        "package_purchase",
		"meeting_sid": "abc123x",
	} {
		if got := params.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestPaymentFilterParams_OmitsUnsetFilters(t *testing.T) {
	cmd := &cobra.Command{Use: "list"}
	addPaymentFilterFlags(cmd, "refunds", "")

	params := paymentFilterParams(cmd)

	if len(params) != 0 {
		t.Errorf("params = %v, want empty", params)
	}
}
