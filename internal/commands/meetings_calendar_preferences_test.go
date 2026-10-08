package commands

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newCalendarPreferencesUpdateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "update"}
	cmd.Flags().Bool("override-calendars", false, "")
	cmd.Flags().String("event-add-calendar-ids", "", "")
	cmd.Flags().Bool("override-conflict-check-calendars", false, "")
	cmd.Flags().String("conflict-check-calendar-ids", "", "")
	cmd.Flags().String("busy-statuses", "", "")
	cmd.Flags().Bool("override-event-layout", false, "")
	cmd.Flags().String("summary-type", "", "")
	cmd.Flags().String("custom-summary", "", "")
	cmd.Flags().String("body", "", "")
	cmd.Flags().String("event-color", "", "")
	cmd.Flags().String("json-file", "", "")
	return cmd
}

func TestCalendarPreferenceUpdateBody_OmitsFlagsTheCallerNeverSet(t *testing.T) {
	cmd := newCalendarPreferencesUpdateTestCmd()
	if err := cmd.Flags().Set("event-color", "#039BE5"); err != nil {
		t.Fatalf("set event-color: %v", err)
	}

	body, err := calendarPreferenceUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	want := map[string]interface{}{"event_color": "#039BE5"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v; unset flags must not reset the other settings", body, want)
	}
}

func TestCalendarPreferenceUpdateBody_SendsEveryFlagTheCallerSet(t *testing.T) {
	cmd := newCalendarPreferencesUpdateTestCmd()
	for flag, value := range map[string]string{
		"override-calendars":                "true",
		"event-add-calendar-ids":            "cal-1, cal-2",
		"override-conflict-check-calendars": "false",
		"conflict-check-calendar-ids":       "cal-3",
		"busy-statuses":                     "busy,oof",
		"override-event-layout":             "true",
		"summary-type":                      "custom",
		"custom-summary":                    "Consultation",
		"body":                              "<p>Agenda</p>",
	} {
		if err := cmd.Flags().Set(flag, value); err != nil {
			t.Fatalf("set %s: %v", flag, err)
		}
	}

	body, err := calendarPreferenceUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	want := map[string]interface{}{
		"override_calendars":                true,
		"event_add_calendar_ids":            []string{"cal-1", "cal-2"},
		"override_conflict_check_calendars": false,
		"conflict_check_calendar_ids":       []string{"cal-3"},
		"busy_statuses":                     []string{"busy", "oof"},
		"override_event_layout":             true,
		"summary_type":                      "custom",
		"custom_summary":                    "Consultation",
		"body":                              "<p>Agenda</p>",
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestCalendarPreferenceUpdateBody_SendsAnEmptyListToClearCalendars(t *testing.T) {
	cmd := newCalendarPreferencesUpdateTestCmd()
	if err := cmd.Flags().Set("conflict-check-calendar-ids", ""); err != nil {
		t.Fatalf("set conflict-check-calendar-ids: %v", err)
	}

	body, err := calendarPreferenceUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	ids, ok := body["conflict_check_calendar_ids"].([]string)
	if !ok || ids == nil || len(ids) != 0 {
		t.Errorf("conflict_check_calendar_ids = %#v, want an empty, non-nil list so the API clears it", body["conflict_check_calendar_ids"])
	}
}

func TestCalendarPreferenceUpdateBody_FlagsOverrideTheJSONFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "preference.json")
	content := `{"summary_type": "meeting_name", "override_event_layout": true}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write json file: %v", err)
	}

	cmd := newCalendarPreferencesUpdateTestCmd()
	for flag, value := range map[string]string{"json-file": path, "summary-type": "custom"} {
		if err := cmd.Flags().Set(flag, value); err != nil {
			t.Fatalf("set %s: %v", flag, err)
		}
	}

	body, err := calendarPreferenceUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	if body["summary_type"] != "custom" {
		t.Errorf("summary_type = %v, want custom from the flag", body["summary_type"])
	}
	if body["override_event_layout"] != true {
		t.Errorf("override_event_layout = %v, want true from the JSON file", body["override_event_layout"])
	}
}
