package commands

import (
	"fmt"

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

		params := paginationParams(cmd)

		hostEmail, _ := cmd.Flags().GetString("host-email")
		if hostEmail != "" {
			params.Set("host_email", hostEmail)
		}

		clientEmail, _ := cmd.Flags().GetString("client-email")
		if clientEmail != "" {
			params.Set("client_email", clientEmail)
		}

		bookingType, _ := cmd.Flags().GetString("type")
		if bookingType != "" {
			params.Set("type", bookingType)
		}

		sortingKey, _ := cmd.Flags().GetString("sorting-key")
		if sortingKey != "" {
			params.Set("sorting_key", sortingKey)
		}

		sortingOrder, _ := cmd.Flags().GetString("sorting-order")
		if sortingOrder != "" {
			params.Set("sorting_order", sortingOrder)
		}

		data, err := c.Get("/bookings", params)
		if err != nil {
			return err
		}

		printList(data, "bookings", nil)
		return nil
	},
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

		body := map[string]interface{}{}

		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			body["name"] = name
		}

		if cmd.Flags().Changed("email") {
			email, _ := cmd.Flags().GetString("email")
			body["email"] = email
		}

		if cmd.Flags().Changed("status") {
			status, _ := cmd.Flags().GetString("status")
			body["status"] = status
		}

		if cmd.Flags().Changed("cancel-reason") {
			cancelReason, _ := cmd.Flags().GetString("cancel-reason")
			body["cancel_reason"] = cancelReason
		}

		if cmd.Flags().Changed("rejection-reason") {
			rejectionReason, _ := cmd.Flags().GetString("rejection-reason")
			body["rejection_reason"] = rejectionReason
		}

		if cmd.Flags().Changed("slot-date") {
			slotDate, _ := cmd.Flags().GetString("slot-date")
			body["slot_date"] = slotDate
		}

		if cmd.Flags().Changed("slot-start-time") {
			slotStartTime, _ := cmd.Flags().GetString("slot-start-time")
			body["slot_start_time"] = slotStartTime
		}

		if cmd.Flags().Changed("time-zone") {
			timeZone, _ := cmd.Flags().GetString("time-zone")
			body["time_zone"] = timeZone
		}

		if cmd.Flags().Changed("reschedule-reason") {
			rescheduleReason, _ := cmd.Flags().GetString("reschedule-reason")
			body["reschedule_reason"] = rescheduleReason
		}

		data, err := c.Put(fmt.Sprintf("/bookings/%s", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
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
}
