package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var availabilitiesCmd = &cobra.Command{
	Use:   "availabilities",
	Short: "Manage availabilities",
}

var availabilitiesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all availabilities",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		params := paginationParams(cmd)

		emails, _ := cmd.Flags().GetString("emails")
		if emails != "" {
			for _, email := range strings.Split(emails, ",") {
				params.Add("emails[]", email)
			}
		}

		data, err := c.Get("/availabilities", params)
		if err != nil {
			return err
		}

		printList(data, "availabilities", nil)
		return nil
	},
}

var availabilitiesShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show an availability",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/availabilities/%s", args[0]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var availabilitiesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an availability",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
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

		data, err := c.Post("/availabilities", body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var availabilitiesUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an availability",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
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

		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			body["name"] = name
		}

		data, err := c.Put(fmt.Sprintf("/availabilities/%s", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(availabilitiesCmd)

	availabilitiesCmd.AddCommand(availabilitiesListCmd)
	addPaginationFlags(availabilitiesListCmd)
	availabilitiesListCmd.Flags().String("emails", "", "Comma-separated emails to filter by")

	availabilitiesCmd.AddCommand(availabilitiesShowCmd)

	availabilitiesCmd.AddCommand(availabilitiesCreateCmd)
	availabilitiesCreateCmd.Flags().String("email", "", "Email address")
	availabilitiesCreateCmd.Flags().String("name", "", "Availability name")
	availabilitiesCreateCmd.Flags().String("time-zone", "", "Time zone")
	availabilitiesCreateCmd.Flags().String("json-file", "", "Path to JSON file with periods/overrides")

	availabilitiesCmd.AddCommand(availabilitiesUpdateCmd)
	availabilitiesUpdateCmd.Flags().String("name", "", "Availability name")
	availabilitiesUpdateCmd.Flags().String("json-file", "", "Path to JSON file with update data")
}
