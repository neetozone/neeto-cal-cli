package commands

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var meetingOutcomesCmd = &cobra.Command{
	Use:   "meeting-outcomes",
	Short: "Manage meeting outcomes",
}

var meetingOutcomesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List meeting outcomes",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/meeting-outcomes", meetingOutcomesListParams(cmd))
		if err != nil {
			return err
		}

		printList(data, "meeting_outcomes", nil)
		return nil
	},
}

func meetingOutcomesListParams(cmd *cobra.Command) url.Values {
	params := url.Values{}

	search, _ := cmd.Flags().GetString("search")
	if search != "" {
		params.Set("search_term", search)
	}

	return params
}

var meetingOutcomesShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a meeting outcome",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/meeting-outcomes/%s", args[0]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingOutcomesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a meeting outcome",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Post("/meeting-outcomes", meetingOutcomeBody(cmd))
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var meetingOutcomesUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a meeting outcome",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Put(fmt.Sprintf("/meeting-outcomes/%s", args[0]), meetingOutcomeBody(cmd))
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func meetingOutcomeBody(cmd *cobra.Command) map[string]interface{} {
	outcome := map[string]interface{}{}

	for _, flag := range []string{"name", "color"} {
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetString(flag)
			outcome[flag] = value
		}
	}

	return map[string]interface{}{"meeting_outcome": outcome}
}

var meetingOutcomesDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a meeting outcome",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		err = c.Delete(fmt.Sprintf("/meeting-outcomes/%s", args[0]))
		if err != nil {
			return err
		}

		printMessage("Meeting outcome deleted.")
		return nil
	},
}

const meetingOutcomeColorUsage = "Color: blue, orange, gray, red, yellow, violet, pink, teal, cyan or green"

func init() {
	register(func(root *cobra.Command) { root.AddCommand(meetingOutcomesCmd) })

	meetingOutcomesCmd.AddCommand(meetingOutcomesListCmd)
	meetingOutcomesListCmd.Flags().String("search", "", "Search meeting outcomes by name")

	meetingOutcomesCmd.AddCommand(meetingOutcomesShowCmd)

	meetingOutcomesCmd.AddCommand(meetingOutcomesCreateCmd)
	meetingOutcomesCreateCmd.Flags().String("name", "", "Outcome name")
	meetingOutcomesCreateCmd.Flags().String("color", "", meetingOutcomeColorUsage)
	markFlagsRequired(meetingOutcomesCreateCmd, "name", "color")

	meetingOutcomesCmd.AddCommand(meetingOutcomesUpdateCmd)
	meetingOutcomesUpdateCmd.Flags().String("name", "", "Outcome name")
	meetingOutcomesUpdateCmd.Flags().String("color", "", meetingOutcomeColorUsage)

	meetingOutcomesCmd.AddCommand(meetingOutcomesDeleteCmd)
}
