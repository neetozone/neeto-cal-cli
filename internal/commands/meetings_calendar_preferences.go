package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var meetingsCalendarPreferencesCmd = &cobra.Command{
	Use:   "calendar-preferences",
	Short: "Manage a meeting's calendar settings",
}

var meetingsCalendarPreferencesShowCmd = &cobra.Command{
	Use:   "show <meeting-sid> <integration>",
	Short: "Show a meeting's calendar settings for google_calendar, outlook or icloud",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/meetings/%s/calendar-preferences/%s", args[0], args[1]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingsCalendarPreferencesUpdateCmd = &cobra.Command{
	Use:   "update <meeting-sid> <integration>",
	Short: "Update a meeting's calendar settings for google_calendar, outlook or icloud",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		body, err := calendarPreferenceUpdateBody(cmd)
		if err != nil {
			return err
		}

		data, err := c.Patch(fmt.Sprintf("/meetings/%s/calendar-preferences/%s", args[0], args[1]), body)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

func calendarPreferenceUpdateBody(cmd *cobra.Command) (map[string]interface{}, error) {
	body := map[string]interface{}{}

	jsonFile, _ := cmd.Flags().GetString("json-file")
	if jsonFile != "" {
		fileData, err := readJSONFile(jsonFile)
		if err != nil {
			return nil, err
		}
		body = fileData
	}

	for flag, field := range map[string]string{
		"override-calendars":                "override_calendars",
		"override-conflict-check-calendars": "override_conflict_check_calendars",
		"override-event-layout":             "override_event_layout",
	} {
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetBool(flag)
			body[field] = value
		}
	}

	for flag, field := range map[string]string{
		"summary-type":   "summary_type",
		"custom-summary": "custom_summary",
		"body":           "body",
		"event-color":    "event_color",
	} {
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetString(flag)
			body[field] = value
		}
	}

	for flag, field := range map[string]string{
		"event-add-calendar-ids":      "event_add_calendar_ids",
		"conflict-check-calendar-ids": "conflict_check_calendar_ids",
		"busy-statuses":               "busy_statuses",
	} {
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetString(flag)
			body[field] = splitCSV(value)
		}
	}

	return body, nil
}

func init() {
	meetingsCmd.AddCommand(meetingsCalendarPreferencesCmd)

	meetingsCalendarPreferencesCmd.AddCommand(meetingsCalendarPreferencesShowCmd)

	meetingsCalendarPreferencesCmd.AddCommand(meetingsCalendarPreferencesUpdateCmd)
	flags := meetingsCalendarPreferencesUpdateCmd.Flags()
	flags.Bool("override-calendars", false, "Add booking events to --event-add-calendar-ids instead of the host's default calendars")
	flags.String("event-add-calendar-ids", "", "Comma-separated calendar IDs to add booking events to, in order. Pass \"\" to clear")
	flags.Bool("override-conflict-check-calendars", false, "Check --conflict-check-calendar-ids for conflicts instead of the host's default calendars")
	flags.String("conflict-check-calendar-ids", "", "Comma-separated calendar IDs to check for conflicts. Pass \"\" to clear")
	flags.String("busy-statuses", "", "Outlook only. Comma-separated event statuses that block a slot: tentative, busy, oof, workingElsewhere")
	flags.Bool("override-event-layout", false, "Use this meeting's own event title and description")
	flags.String("summary-type", "", "Event title format: host_and_client, client_and_host, meeting_name or custom")
	flags.String("custom-summary", "", "Event title, used when --summary-type is custom")
	flags.String("body", "", "Event description (HTML). Pass \"\" to reset it to the default")
	flags.String("event-color", "", "Google Calendar only. Event color hex, e.g. #039BE5. Pass \"\" to reset it")
	flags.String("json-file", "", "Path to a JSON file with the calendar settings")
}
