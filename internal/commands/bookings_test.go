package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newBookingsListTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "list"}
	addPaginationFlags(cmd)
	for _, flag := range []string{
		"host-email", "client-email", "type", "sorting-key", "sorting-order",
		"search", "meeting-sid", "starts-after", "starts-before",
	} {
		cmd.Flags().String(flag, "", "")
	}
	return cmd
}

func newBookingsUpdateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "update"}
	for _, flag := range []string{
		"name", "email", "status", "cancel-reason", "rejection-reason", "slot-date",
		"slot-start-time", "time-zone", "reschedule-reason", "preferred-meeting-spot", "meeting-outcome-id",
	} {
		cmd.Flags().String(flag, "", "")
	}
	cmd.Flags().Bool("override-availability", false, "")
	return cmd
}

func newBookingsCreateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create"}
	for _, flag := range []string{
		"meeting-slug", "email", "name", "slot-date", "slot-start-time", "time-zone",
		"preferred-meeting-spot", "host-email",
	} {
		cmd.Flags().String(flag, "", "")
	}
	cmd.Flags().Bool("override-availability", false, "")
	return cmd
}

func setFlags(t *testing.T, cmd *cobra.Command, values map[string]string) {
	t.Helper()
	for flag, value := range values {
		if err := cmd.Flags().Set(flag, value); err != nil {
			t.Fatalf("set %s: %v", flag, err)
		}
	}
}

var requiredBookingCreateFlags = map[string]string{
	"meeting-slug":    "meeting-with-oliver",
	"email":           "eve@example.com",
	"name":            "Eve Smith",
	"slot-date":       "2026-10-05",
	"slot-start-time": "08:00 PM",
	"time-zone":       "Asia/Kolkata",
}

var requiredBookingCreateBody = map[string]interface{}{
	"meeting_slug":    "meeting-with-oliver",
	"email":           "eve@example.com",
	"name":            "Eve Smith",
	"slot_date":       "2026-10-05",
	"slot_start_time": "08:00 PM",
	"time_zone":       "Asia/Kolkata",
}

func TestBookingCreateBody_SendsOnlyTheRequiredFieldsByDefault(t *testing.T) {
	cmd := newBookingsCreateTestCmd()
	setFlags(t, cmd, requiredBookingCreateFlags)

	body := bookingCreateBody(cmd)

	if !reflect.DeepEqual(body, requiredBookingCreateBody) {
		t.Errorf("body = %v, want %v", body, requiredBookingCreateBody)
	}
}

func TestBookingCreateBody_SendsOverrideAvailabilityAndHostEmail(t *testing.T) {
	cmd := newBookingsCreateTestCmd()
	setFlags(t, cmd, requiredBookingCreateFlags)
	setFlags(t, cmd, map[string]string{
		"override-availability":  "true",
		"host-email":             "oliver@example.com",
		"preferred-meeting-spot": "zoom",
	})

	body := bookingCreateBody(cmd)

	want := map[string]interface{}{
		"override_availability":  true,
		"host_email":             "oliver@example.com",
		"preferred_meeting_spot": "zoom",
	}
	for key, value := range requiredBookingCreateBody {
		want[key] = value
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestBookingUpdateBody_SendsOverrideAvailabilityWithTheReschedule(t *testing.T) {
	cmd := newBookingsUpdateTestCmd()
	setFlags(t, cmd, map[string]string{
		"slot-date":             "2026-10-05",
		"slot-start-time":       "09:00 PM",
		"time-zone":             "Asia/Kolkata",
		"override-availability": "true",
	})

	body := bookingUpdateBody(cmd)

	want := map[string]interface{}{
		"slot_date":             "2026-10-05",
		"slot_start_time":       "09:00 PM",
		"time_zone":             "Asia/Kolkata",
		"override_availability": true,
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestBookingsListParams_SendsSearchAndDateRangeFilters(t *testing.T) {
	cmd := newBookingsListTestCmd()
	for flag, value := range map[string]string{
		"search":        "acme.com",
		"meeting-sid":   "abc123x",
		"starts-after":  "2026-10-01",
		"starts-before": "2026-10-31T23:59:59-04:00",
	} {
		if err := cmd.Flags().Set(flag, value); err != nil {
			t.Fatalf("set %s: %v", flag, err)
		}
	}

	params := bookingsListParams(cmd)

	for param, want := range map[string]string{
		"search":        "acme.com",
		"meeting_sid":   "abc123x",
		"starts_after":  "2026-10-01",
		"starts_before": "2026-10-31T23:59:59-04:00",
	} {
		if got := params.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestBookingsListParams_Empty(t *testing.T) {
	cmd := newBookingsListTestCmd()

	params := bookingsListParams(cmd)

	if len(params) != 0 {
		t.Errorf("params = %v, want empty", params)
	}
}

func TestBookingUpdateBody_SendsRescheduleSpotAndOutcome(t *testing.T) {
	cmd := newBookingsUpdateTestCmd()
	for flag, value := range map[string]string{
		"slot-date":              "2026-10-06",
		"slot-start-time":        "10:00",
		"preferred-meeting-spot": "zoom",
		"meeting-outcome-id":     "outcome-1",
	} {
		if err := cmd.Flags().Set(flag, value); err != nil {
			t.Fatalf("set %s: %v", flag, err)
		}
	}

	body := bookingUpdateBody(cmd)

	want := map[string]interface{}{
		"slot_date":              "2026-10-06",
		"slot_start_time":        "10:00",
		"preferred_meeting_spot": "zoom",
		"meeting_outcome_id":     "outcome-1",
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestBookingUpdateBody_SendsAnExplicitlyEmptyOutcomeToClearIt(t *testing.T) {
	cmd := newBookingsUpdateTestCmd()
	if err := cmd.Flags().Set("meeting-outcome-id", ""); err != nil {
		t.Fatalf("set meeting-outcome-id: %v", err)
	}

	body := bookingUpdateBody(cmd)

	want := map[string]interface{}{"meeting_outcome_id": ""}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestBookingUpdateBody_OmitsFlagsTheCallerNeverSet(t *testing.T) {
	cmd := newBookingsUpdateTestCmd()

	body := bookingUpdateBody(cmd)

	if len(body) != 0 {
		t.Errorf("body = %v, want empty", body)
	}
}
