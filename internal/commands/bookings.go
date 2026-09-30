package commands

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var bookingsCmd = &cobra.Command{
	Use:   "bookings",
	Short: "Manage bookings",
}

var bookingsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List bookings",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		params := bookingsListParams(cmd)

		data, err := c.Get("/bookings", params)
		if err != nil {
			return err
		}

		printList(data, "bookings", nil)
		return nil
	},
}

func bookingsListParams(cmd *cobra.Command) url.Values {
	params := paginationParams(cmd)

	for flag, param := range map[string]string{
		"host-email":    "host_email",
		"client-email":  "client_email",
		"type":          "type",
		"sorting-key":   "sorting_key",
		"sorting-order": "sorting_order",
		"search":        "search",
		"meeting-sid":   "meeting_sid",
		"starts-after":  "starts_after",
		"starts-before": "starts_before",
	} {
		value, _ := cmd.Flags().GetString(flag)
		if value != "" {
			params.Set(param, value)
		}
	}

	return params
}

var bookingsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a booking",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/bookings/%s", args[0]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var bookingsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a booking",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		meetingSlug, _ := cmd.Flags().GetString("meeting-slug")
		email, _ := cmd.Flags().GetString("email")
		name, _ := cmd.Flags().GetString("name")
		slotDate, _ := cmd.Flags().GetString("slot-date")
		slotStartTime, _ := cmd.Flags().GetString("slot-start-time")
		timeZone, _ := cmd.Flags().GetString("time-zone")

		body := map[string]interface{}{
			"meeting_slug":    meetingSlug,
			"email":           email,
			"name":            name,
			"slot_date":       slotDate,
			"slot_start_time": slotStartTime,
			"time_zone":       timeZone,
		}

		preferredSpot, _ := cmd.Flags().GetString("preferred-meeting-spot")
		if preferredSpot != "" {
			body["preferred_meeting_spot"] = preferredSpot
		}

		data, err := c.Post("/bookings", body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var bookingsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a booking",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		body := bookingUpdateBody(cmd)

		data, err := c.Put(fmt.Sprintf("/bookings/%s", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func bookingUpdateBody(cmd *cobra.Command) map[string]interface{} {
	body := map[string]interface{}{}

	for flag, key := range map[string]string{
		"name":                   "name",
		"email":                  "email",
		"status":                 "status",
		"cancel-reason":          "cancel_reason",
		"rejection-reason":       "rejection_reason",
		"slot-date":              "slot_date",
		"slot-start-time":        "slot_start_time",
		"time-zone":              "time_zone",
		"reschedule-reason":      "reschedule_reason",
		"preferred-meeting-spot": "preferred_meeting_spot",
		"meeting-outcome-id":     "meeting_outcome_id",
	} {
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetString(flag)
			body[key] = value
		}
	}

	return body
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(bookingsCmd) })

	bookingsCmd.AddCommand(bookingsListCmd)
	addPaginationFlags(bookingsListCmd)
	bookingsListCmd.Flags().String("host-email", "", "Filter by host email")
	bookingsListCmd.Flags().String("client-email", "", "Filter by client email")
	bookingsListCmd.Flags().String("type", "", "Filter by type (upcoming, past, cancelled, incomplete)")
	bookingsListCmd.Flags().String("sorting-key", "", "Sort by field (created_at, starts_at)")
	bookingsListCmd.Flags().String("sorting-order", "", "Sort order (asc, desc)")
	bookingsListCmd.Flags().String("search", "", "Search text, min 3 chars (client name/email/sid, meeting name, host name/email, form answers)")
	bookingsListCmd.Flags().String("meeting-sid", "", "Filter by meeting SID")
	bookingsListCmd.Flags().String("starts-after", "", "Only bookings starting at or after this time (ISO 8601 datetime or YYYY-MM-DD)")
	bookingsListCmd.Flags().String("starts-before", "", "Only bookings starting before this time (ISO 8601 datetime or YYYY-MM-DD)")

	bookingsCmd.AddCommand(bookingsShowCmd)

	bookingsCmd.AddCommand(bookingsCreateCmd)
	bookingsCreateCmd.Flags().String("meeting-slug", "", "Meeting slug")
	bookingsCreateCmd.Flags().String("email", "", "Client email")
	bookingsCreateCmd.Flags().String("name", "", "Client name")
	bookingsCreateCmd.Flags().String("slot-date", "", "Slot date (YYYY-MM-DD)")
	bookingsCreateCmd.Flags().String("slot-start-time", "", "Slot start time (HH:MM)")
	bookingsCreateCmd.Flags().String("time-zone", "", timeZoneUsage)
	bookingsCreateCmd.Flags().String("preferred-meeting-spot", "", "Preferred meeting spot")
	markFlagsRequired(bookingsCreateCmd, "meeting-slug", "email", "name", "slot-date", "slot-start-time", "time-zone")

	bookingsCmd.AddCommand(bookingsUpdateCmd)
	bookingsUpdateCmd.Flags().String("name", "", "Client name (optional; overrides the existing booking's name when rescheduling)")
	bookingsUpdateCmd.Flags().String("email", "", "Client email (optional; overrides the existing booking's email when rescheduling)")
	bookingsUpdateCmd.Flags().String("status", "", "Status (cancelled, approved, rejected)")
	bookingsUpdateCmd.Flags().String("cancel-reason", "", "Cancellation reason")
	bookingsUpdateCmd.Flags().String("rejection-reason", "", "Rejection reason")
	bookingsUpdateCmd.Flags().String("slot-date", "", "New slot date (YYYY-MM-DD)")
	bookingsUpdateCmd.Flags().String("slot-start-time", "", "New slot start time (HH:MM)")
	bookingsUpdateCmd.Flags().String("time-zone", "", timeZoneUsage)
	bookingsUpdateCmd.Flags().String("reschedule-reason", "", "Reschedule reason")
	bookingsUpdateCmd.Flags().String("preferred-meeting-spot", "", "Preferred meeting spot for the rescheduled slot (use with --slot-date and --slot-start-time)")
	bookingsUpdateCmd.Flags().String("meeting-outcome-id", "", "Meeting outcome ID (pass an empty string to clear the outcome)")
}
