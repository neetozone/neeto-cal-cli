package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func newMeetingsUnavailableSlotsTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "unavailable-slots"}
	cmd.Flags().Int("year", 0, "")
	cmd.Flags().Int("month", 0, "")
	cmd.Flags().Int("day", 0, "")
	cmd.Flags().String("time-zone", "", "")
	cmd.Flags().String("host-email", "", "")
	cmd.Flags().Int("duration", 0, "")
	return cmd
}

func TestMeetingUnavailableSlotsParams_SendsEveryFlag(t *testing.T) {
	cmd := newMeetingsUnavailableSlotsTestCmd()
	setFlags(t, cmd, map[string]string{
		"year": "2026", "month": "10", "day": "19", "time-zone": "America/New_York",
		"host-email": "oliver@example.com", "duration": "60",
	})

	params := meetingUnavailableSlotsParams(cmd)

	for param, want := range map[string]string{
		"year": "2026", "month": "10", "day": "19", "time_zone": "America/New_York",
		"host_email": "oliver@example.com", "duration": "60",
	} {
		if got := params.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestMeetingUnavailableSlotsParams_OmitsOptionalFlagsWhenUnset(t *testing.T) {
	cmd := newMeetingsUnavailableSlotsTestCmd()
	setFlags(t, cmd, map[string]string{"year": "2026", "month": "10", "time-zone": "America/New_York"})

	params := meetingUnavailableSlotsParams(cmd)

	for _, param := range []string{"day", "host_email", "duration"} {
		if params.Has(param) {
			t.Errorf("%s = %q, want absent", param, params.Get(param))
		}
	}
}
