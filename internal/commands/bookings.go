package commands

import (
	"errors"
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
		"host-email":        "host_email",
		"client-email":      "client_email",
		"type":              "type",
		"sorting-key":       "sorting_key",
		"sorting-order":     "sorting_order",
		"search":            "search",
		"meeting-sid":       "meeting_sid",
		"starts-after":      "starts_after",
		"starts-before":     "starts_before",
		"parent-booking-id": "parent_booking_id",
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

		body, err := bookingCreateBody(cmd)
		if err != nil {
			return err
		}

		data, err := c.Post("/bookings", body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func bookingCreateBody(cmd *cobra.Command) (map[string]interface{}, error) {
	body := map[string]interface{}{}

	for flag, key := range map[string]string{
		"meeting-slug":    "meeting_slug",
		"name":            "name",
		"slot-date":       "slot_date",
		"slot-start-time": "slot_start_time",
		"time-zone":       "time_zone",
	} {
		value, _ := cmd.Flags().GetString(flag)
		body[key] = value
	}

	email, _ := cmd.Flags().GetString("email")
	phoneNumber, _ := cmd.Flags().GetString("phonenumber")
	if email == "" && phoneNumber == "" {
		return nil, errors.New("pass --email or --phonenumber")
	}

	for flag, key := range map[string]string{
		"preferred-meeting-spot": "preferred_meeting_spot",
		"host-email":             "host_email",
		"email":                  "email",
		"phonenumber":            "phonenumber",
	} {
		value, _ := cmd.Flags().GetString(flag)
		if value != "" {
			body[key] = value
		}
	}

	if err := setJSONObjectFlag(cmd, body, "metadata", "metadata"); err != nil {
		return nil, err
	}

	if err := setJSONObjectFlag(cmd, body, "form-responses", "form_responses"); err != nil {
		return nil, err
	}

	addOnServiceIDs, _ := cmd.Flags().GetString("add-on-service-ids")
	if addOnServiceIDs != "" {
		body["add_on_service_ids"] = splitCSV(addOnServiceIDs)
	}

	if reserveSlot, _ := cmd.Flags().GetBool("reserve-slot"); reserveSlot {
		body["reserve_slot"] = true
	}

	if cmd.Flags().Changed("send-payment-link") {
		sendPaymentLink, _ := cmd.Flags().GetBool("send-payment-link")
		body["send_payment_link"] = sendPaymentLink
	}

	if override, _ := cmd.Flags().GetBool("override-availability"); override {
		body["override_availability"] = true
	}

	return body, nil
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

		body, err := bookingUpdateBody(cmd)
		if err != nil {
			return err
		}

		data, err := c.Put(fmt.Sprintf("/bookings/%s", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func bookingUpdateBody(cmd *cobra.Command) (map[string]interface{}, error) {
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
		"phonenumber":            "phonenumber",
	} {
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetString(flag)
			body[key] = value
		}
	}

	if cmd.Flags().Changed("metadata") {
		if err := setJSONObjectFlag(cmd, body, "metadata", "metadata"); err != nil {
			return nil, err
		}
	}

	if override, _ := cmd.Flags().GetBool("override-availability"); override {
		body["override_availability"] = true
	}

	return body, nil
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
	bookingsListCmd.Flags().String("parent-booking-id", "", "Only bookings created together with this booking, such as pre-bookings")

	bookingsCmd.AddCommand(bookingsShowCmd)

	bookingsCmd.AddCommand(bookingsCreateCmd)
	bookingsCreateCmd.Flags().String("meeting-slug", "", "Meeting slug")
	bookingsCreateCmd.Flags().String("email", "", "Client email (required unless the meeting identifies clients by phone)")
	bookingsCreateCmd.Flags().String("name", "", "Client name")
	bookingsCreateCmd.Flags().String("slot-date", "", "Slot date (YYYY-MM-DD)")
	bookingsCreateCmd.Flags().String("slot-start-time", "", "Slot start time (HH:MM)")
	bookingsCreateCmd.Flags().String("time-zone", "", timeZoneUsage)
	bookingsCreateCmd.Flags().String("preferred-meeting-spot", "", "Preferred meeting spot")
	bookingsCreateCmd.Flags().Bool("override-availability", false, "Book the slot even if it is outside the meeting's availability (host or admin only)")
	bookingsCreateCmd.Flags().String("host-email", "", "Host to assign on a meeting with several hosts (needs --override-availability unless clients can choose the host)")
	bookingsCreateCmd.Flags().String("phonenumber", "", "Client phone number in international format (required when the meeting identifies clients by phone)")
	bookingsCreateCmd.Flags().String("metadata", "", "Booking metadata as a JSON object, such as '{\"source\":\"crm\"}'")
	bookingsCreateCmd.Flags().String("form-responses", "", "Booking form answers as a JSON object keyed by question field code")
	bookingsCreateCmd.Flags().String("add-on-service-ids", "", "Comma-separated IDs of the meeting's add-on services")
	bookingsCreateCmd.Flags().Bool("reserve-slot", false, "Reserve the slot without sending a payment link")
	bookingsCreateCmd.Flags().Bool("send-payment-link", true, "Send a payment link for a paid meeting (pass --send-payment-link=false to skip)")
	markFlagsRequired(bookingsCreateCmd, "meeting-slug", "name", "slot-date", "slot-start-time", "time-zone")

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
	bookingsUpdateCmd.Flags().String("phonenumber", "", "Client phone number (optional; overrides the existing booking's phone number when rescheduling)")
	bookingsUpdateCmd.Flags().String("metadata", "", "Booking metadata as a JSON object, used when rescheduling")
	bookingsUpdateCmd.Flags().Bool("override-availability", false, "Reschedule even if the new slot is outside the meeting's availability (host or admin only; use with --slot-date and --slot-start-time)")
}
