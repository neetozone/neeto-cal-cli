package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var meetingsAvailabilitiesCmd = &cobra.Command{
	Use:   "availabilities",
	Short: "Manage meeting-level availabilities",
}

var meetingsAvailabilitiesCreateCmd = &cobra.Command{
	Use:   "create <meeting-sid>",
	Short: "Create availability for a meeting",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{}

		jsonFile, _ := cmd.Flags().GetString("json-file")
		if jsonFile != "" {
			fileData, err := readJSONFile(jsonFile)
			if err != nil {
				return err
			}
			body = fileData
		}

		email, _ := cmd.Flags().GetString("email")
		if email != "" {
			body["email"] = email
		}

		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			body["name"] = name
		}

		timeZone, _ := cmd.Flags().GetString("time-zone")
		if timeZone != "" {
			body["time_zone"] = timeZone
		}

		data, err := c.Post(fmt.Sprintf("/meetings/%s/availabilities", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var meetingsAvailabilitiesUpdateCmd = &cobra.Command{
	Use:   "update <meeting-sid>",
	Short: "Update availability for a meeting",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{}

		jsonFile, _ := cmd.Flags().GetString("json-file")
		if jsonFile != "" {
			fileData, err := readJSONFile(jsonFile)
			if err != nil {
				return err
			}
			body = fileData
		}

		data, err := c.Patch(fmt.Sprintf("/meetings/%s/availabilities", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func init() {
	meetingsCmd.AddCommand(meetingsAvailabilitiesCmd)

	meetingsAvailabilitiesCmd.AddCommand(meetingsAvailabilitiesCreateCmd)
	meetingsAvailabilitiesCreateCmd.Flags().String("email", "", "Email address")
	meetingsAvailabilitiesCreateCmd.Flags().String("name", "", "Availability name")
	meetingsAvailabilitiesCreateCmd.Flags().String("time-zone", "", "Time zone")
	meetingsAvailabilitiesCreateCmd.Flags().String("json-file", "", "Path to JSON file with periods/overrides")

	meetingsAvailabilitiesCmd.AddCommand(meetingsAvailabilitiesUpdateCmd)
	meetingsAvailabilitiesUpdateCmd.Flags().String("json-file", "", "Path to JSON file with update data")
}
