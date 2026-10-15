package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func newHolidaysListTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "list"}
	addPaginationFlags(cmd)
	for _, flag := range []string{"search", "from-date", "to-date"} {
		cmd.Flags().String(flag, "", "")
	}
	cmd.Flags().Bool("include-disabled", false, "")
	return cmd
}

func TestHolidaysListParams_SendsEveryFilter(t *testing.T) {
	cmd := newHolidaysListTestCmd()
	setFlags(t, cmd, map[string]string{
		"search": "offsite", "from-date": "2026-10-01", "to-date": "2026-10-31", "include-disabled": "true",
	})

	params := holidaysListParams(cmd)

	for param, want := range map[string]string{
		"search": "offsite", "from_date": "2026-10-01", "to_date": "2026-10-31", "include_disabled": "true",
	} {
		if got := params.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestHolidaysListParams_Empty(t *testing.T) {
	params := holidaysListParams(newHolidaysListTestCmd())

	if len(params) != 0 {
		t.Errorf("params = %v, want empty", params)
	}
}
