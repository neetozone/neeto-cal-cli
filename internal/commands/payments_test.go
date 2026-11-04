package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func newListFilterTestCmd(flags ...string) *cobra.Command {
	cmd := &cobra.Command{Use: "list"}
	for _, flag := range flags {
		cmd.Flags().String(flag, "", "")
	}
	addPaginationFlags(cmd)
	return cmd
}

func TestPaymentsListParams_SendsFilters(t *testing.T) {
	cmd := newListFilterTestCmd("paid-after", "paid-before", "status", "type", "meeting-sid")
	setFlags(t, cmd, map[string]string{
		"paid-after":  "2026-10-01",
		"paid-before": "2026-11-01",
		"status":      "declined",
		"type":        "package_purchase",
		"meeting-sid": "abc123x",
	})

	params := paymentsListParams(cmd)

	for param, want := range map[string]string{
		"paid_after":  "2026-10-01",
		"paid_before": "2026-11-01",
		"status":      "declined",
		"type":        "package_purchase",
		"meeting_sid": "abc123x",
	} {
		if got := params.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestRefundsListParams_SendsFilters(t *testing.T) {
	cmd := newListFilterTestCmd("refunded-after", "refunded-before", "status", "type", "meeting-sid")
	setFlags(t, cmd, map[string]string{"refunded-after": "2026-10-01", "refunded-before": "2026-11-01"})

	params := refundsListParams(cmd)

	for param, want := range map[string]string{"refunded_after": "2026-10-01", "refunded_before": "2026-11-01"} {
		if got := params.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
	if params.Has("status") {
		t.Errorf("status = %q, want unset", params.Get("status"))
	}
}
