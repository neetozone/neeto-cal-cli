package commands

import (
	"fmt"

	"github.com/neetozone/neeto-cal-cli/internal/output"
	"github.com/spf13/cobra"
)

var meetingsSpotsCmd = &cobra.Command{
	Use:   "spots",
	Short: "Manage meeting spots",
}

var meetingsSpotsListCmd = &cobra.Command{
	Use:   "list <meeting-sid>",
	Short: "List spots for a meeting",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/meetings/%s/spots", args[0]), nil)
		if err != nil {
			return err
		}

		printList(data, "spots", nil)
		return nil
	},
}

var meetingsSpotsShowCmd = &cobra.Command{
	Use:   "show <meeting-sid> <id>",
	Short: "Show a meeting spot",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/meetings/%s/spots/%s", args[0], args[1]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingsSpotsCreateCmd = &cobra.Command{
	Use:   "create <meeting-sid>",
	Short: "Create a meeting spot",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		spot, _ := cmd.Flags().GetString("spot")
		isDefault, _ := cmd.Flags().GetBool("is-default")

		body := map[string]interface{}{
			"spot":       spot,
			"is_default": isDefault,
		}

		phoneNumber, _ := cmd.Flags().GetString("phone-number")
		if phoneNumber != "" {
			body["phone_number"] = phoneNumber
		}

		location, _ := cmd.Flags().GetString("location")
		if location != "" {
			body["location"] = location
		}

		customText, _ := cmd.Flags().GetString("custom-text")
		if customText != "" {
			body["custom_text"] = customText
		}

		data, err := c.Post(fmt.Sprintf("/meetings/%s/spots", args[0]), body)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingsSpotsUpdateCmd = &cobra.Command{
	Use:   "update <meeting-sid> <id>",
	Short: "Update a meeting spot",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{}

		if cmd.Flags().Changed("spot") {
			spot, _ := cmd.Flags().GetString("spot")
			body["spot"] = spot
		}

		if cmd.Flags().Changed("is-default") {
			isDefault, _ := cmd.Flags().GetBool("is-default")
			body["is_default"] = isDefault
		}

		if cmd.Flags().Changed("phone-number") {
			phoneNumber, _ := cmd.Flags().GetString("phone-number")
			body["phone_number"] = phoneNumber
		}

		if cmd.Flags().Changed("location") {
			location, _ := cmd.Flags().GetString("location")
			body["location"] = location
		}

		if cmd.Flags().Changed("custom-text") {
			customText, _ := cmd.Flags().GetString("custom-text")
			body["custom_text"] = customText
		}

		data, err := c.Put(fmt.Sprintf("/meetings/%s/spots/%s", args[0], args[1]), body)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingsSpotsDeleteCmd = &cobra.Command{
	Use:   "delete <meeting-sid> <id>",
	Short: "Delete a meeting spot",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		err = c.Delete(fmt.Sprintf("/meetings/%s/spots/%s", args[0], args[1]))
		if err != nil {
			return err
		}

		output.PrintMessage("Spot deleted successfully.")
		return nil
	},
}

func init() {
	meetingsCmd.AddCommand(meetingsSpotsCmd)

	meetingsSpotsCmd.AddCommand(meetingsSpotsListCmd)
	meetingsSpotsCmd.AddCommand(meetingsSpotsShowCmd)

	meetingsSpotsCmd.AddCommand(meetingsSpotsCreateCmd)
	meetingsSpotsCreateCmd.Flags().String("spot", "", "Spot type")
	meetingsSpotsCreateCmd.Flags().Bool("is-default", false, "Set as default spot")
	meetingsSpotsCreateCmd.Flags().String("phone-number", "", "Phone number for phone spots")
	meetingsSpotsCreateCmd.Flags().String("location", "", "Location for in-person spots")
	meetingsSpotsCreateCmd.Flags().String("custom-text", "", "Custom text for the spot")
	_ = meetingsSpotsCreateCmd.MarkFlagRequired("spot")

	meetingsSpotsCmd.AddCommand(meetingsSpotsUpdateCmd)
	meetingsSpotsUpdateCmd.Flags().String("spot", "", "Spot type")
	meetingsSpotsUpdateCmd.Flags().Bool("is-default", false, "Set as default spot")
	meetingsSpotsUpdateCmd.Flags().String("phone-number", "", "Phone number for phone spots")
	meetingsSpotsUpdateCmd.Flags().String("location", "", "Location for in-person spots")
	meetingsSpotsUpdateCmd.Flags().String("custom-text", "", "Custom text for the spot")

	meetingsSpotsCmd.AddCommand(meetingsSpotsDeleteCmd)
}
