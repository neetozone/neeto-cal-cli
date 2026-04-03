package commands

import (
	"fmt"

	"github.com/neetozone/neeto-cal-cli/internal/output"
	"github.com/spf13/cobra"
)

var meetingsDurationsCmd = &cobra.Command{
	Use:   "durations",
	Short: "Manage meeting durations",
}

var meetingsDurationsListCmd = &cobra.Command{
	Use:   "list <meeting-sid>",
	Short: "List durations for a meeting",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/meetings/%s/durations", args[0]), nil)
		if err != nil {
			return err
		}

		printList(data, "durations", nil)
		return nil
	},
}

var meetingsDurationsShowCmd = &cobra.Command{
	Use:   "show <meeting-sid> <id>",
	Short: "Show a meeting duration",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/meetings/%s/durations/%s", args[0], args[1]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingsDurationsCreateCmd = &cobra.Command{
	Use:   "create <meeting-sid>",
	Short: "Create a meeting duration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		duration, _ := cmd.Flags().GetInt("duration")
		isDefault, _ := cmd.Flags().GetBool("is-default")

		body := map[string]interface{}{
			"duration":   duration,
			"is_default": isDefault,
		}

		data, err := c.Post(fmt.Sprintf("/meetings/%s/durations", args[0]), body)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingsDurationsUpdateCmd = &cobra.Command{
	Use:   "update <meeting-sid> <id>",
	Short: "Update a meeting duration",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{}

		if cmd.Flags().Changed("duration") {
			duration, _ := cmd.Flags().GetInt("duration")
			body["duration"] = duration
		}

		if cmd.Flags().Changed("is-default") {
			isDefault, _ := cmd.Flags().GetBool("is-default")
			body["is_default"] = isDefault
		}

		data, err := c.Put(fmt.Sprintf("/meetings/%s/durations/%s", args[0], args[1]), body)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingsDurationsDeleteCmd = &cobra.Command{
	Use:   "delete <meeting-sid> <id>",
	Short: "Delete a meeting duration",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		err = c.Delete(fmt.Sprintf("/meetings/%s/durations/%s", args[0], args[1]))
		if err != nil {
			return err
		}

		output.PrintMessage("Duration deleted successfully.")
		return nil
	},
}

func init() {
	meetingsCmd.AddCommand(meetingsDurationsCmd)

	meetingsDurationsCmd.AddCommand(meetingsDurationsListCmd)
	meetingsDurationsCmd.AddCommand(meetingsDurationsShowCmd)

	meetingsDurationsCmd.AddCommand(meetingsDurationsCreateCmd)
	meetingsDurationsCreateCmd.Flags().Int("duration", 0, "Duration in minutes")
	meetingsDurationsCreateCmd.Flags().Bool("is-default", false, "Set as default duration")
	_ = meetingsDurationsCreateCmd.MarkFlagRequired("duration")

	meetingsDurationsCmd.AddCommand(meetingsDurationsUpdateCmd)
	meetingsDurationsUpdateCmd.Flags().Int("duration", 0, "Duration in minutes")
	meetingsDurationsUpdateCmd.Flags().Bool("is-default", false, "Set as default duration")

	meetingsDurationsCmd.AddCommand(meetingsDurationsDeleteCmd)
}
