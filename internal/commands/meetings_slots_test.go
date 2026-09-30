package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func newMeetingsSlotsTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "slots"}
	cmd.Flags().Int("year", 0, "")
	cmd.Flags().Int("month", 0, "")
	cmd.Flags().Int("day", 0, "")
	cmd.Flags().String("time-zone", "", "")
	cmd.Flags().String("host-email", "", "")
	return cmd
}

func TestMeetingSlotsParams_SendsHostEmail(t *testing.T) {
	cmd := newMeetingsSlotsTestCmd()
	for flag, value := range map[string]string{
		"year":       "2026",
		"month":      "10",
		"time-zone":  "America/New_York",
		"host-email": "oliver@example.com",
	} {
		if err := cmd.Flags().Set(flag, value); err != nil {
			t.Fatalf("set %s: %v", flag, err)
		}
	}

	params := meetingSlotsParams(cmd)

	if got := params.Get("host_email"); got != "oliver@example.com" {
		t.Errorf("host_email = %q, want oliver@example.com", got)
	}
	if params.Has("day") {
		t.Errorf("day = %q, want absent when --day is not set", params.Get("day"))
	}
}

func TestMeetingSlotsParams_OmitsHostEmailWhenUnset(t *testing.T) {
	cmd := newMeetingsSlotsTestCmd()

	params := meetingSlotsParams(cmd)

	if params.Has("host_email") {
		t.Errorf("host_email = %q, want absent", params.Get("host_email"))
	}
}
