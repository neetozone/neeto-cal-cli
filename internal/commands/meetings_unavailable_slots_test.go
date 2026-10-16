package commands

import (
	"net/url"
	"reflect"
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

	want := url.Values{
		"year": {"2026"}, "month": {"10"}, "day": {"19"}, "time_zone": {"America/New_York"},
		"host_email": {"oliver@example.com"}, "duration": {"60"},
	}
	if !reflect.DeepEqual(params, want) {
		t.Errorf("params = %v, want %v", params, want)
	}
}

func TestMeetingUnavailableSlotsParams_OmitsOptionalFlagsWhenUnset(t *testing.T) {
	cmd := newMeetingsUnavailableSlotsTestCmd()
	setFlags(t, cmd, map[string]string{"year": "2026", "month": "10", "time-zone": "America/New_York"})

	params := meetingUnavailableSlotsParams(cmd)

	want := url.Values{"year": {"2026"}, "month": {"10"}, "time_zone": {"America/New_York"}}
	if !reflect.DeepEqual(params, want) {
		t.Errorf("params = %v, want %v", params, want)
	}
}
