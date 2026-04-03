package commands

import (
	"fmt"
	"strings"

	"github.com/bigbinary/neeto-cal-cli/internal/output"
	"github.com/spf13/cobra"
)

var meetingTemplatesCmd = &cobra.Command{
	Use:   "meeting-templates",
	Short: "Manage meeting templates",
}

var meetingTemplatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all meeting templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		params := paginationParams(cmd)

		hostEmail, _ := cmd.Flags().GetString("host-email")
		if hostEmail != "" {
			params.Set("host_email", hostEmail)
		}

		search, _ := cmd.Flags().GetString("search")
		if search != "" {
			params.Set("search", search)
		}

		data, err := c.Get("/meeting-templates", params)
		if err != nil {
			return err
		}

		printList(data, "meeting_templates", nil)
		return nil
	},
}

var meetingTemplatesShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a meeting template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/meeting-templates/%s", args[0]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingTemplatesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a meeting template",
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

		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			body["name"] = name
		}

		slug, _ := cmd.Flags().GetString("slug")
		if slug != "" {
			body["slug"] = slug
		}

		hosts, _ := cmd.Flags().GetString("hosts")
		if hosts != "" {
			body["hosts"] = strings.Split(hosts, ",")
		}

		kind, _ := cmd.Flags().GetString("kind")
		if kind != "" {
			body["kind"] = kind
		}

		spot, _ := cmd.Flags().GetString("spot")
		if spot != "" {
			body["spot"] = spot
		}

		duration, _ := cmd.Flags().GetInt("duration")
		if duration > 0 {
			body["duration"] = duration
		}

		data, err := c.Post("/meeting-templates", body)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingTemplatesUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a meeting template",
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

		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			body["name"] = name
		}

		slug, _ := cmd.Flags().GetString("slug")
		if slug != "" {
			body["slug"] = slug
		}

		data, err := c.Put(fmt.Sprintf("/meeting-templates/%s", args[0]), body)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingTemplatesDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a meeting template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		err = c.Delete(fmt.Sprintf("/meeting-templates/%s", args[0]))
		if err != nil {
			return err
		}

		output.PrintMessage("Meeting template deleted successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(meetingTemplatesCmd)

	meetingTemplatesCmd.AddCommand(meetingTemplatesListCmd)
	addPaginationFlags(meetingTemplatesListCmd)
	meetingTemplatesListCmd.Flags().String("host-email", "", "Filter by host email")
	meetingTemplatesListCmd.Flags().String("search", "", "Search meeting templates")

	meetingTemplatesCmd.AddCommand(meetingTemplatesShowCmd)

	meetingTemplatesCmd.AddCommand(meetingTemplatesCreateCmd)
	meetingTemplatesCreateCmd.Flags().String("name", "", "Template name")
	meetingTemplatesCreateCmd.Flags().String("slug", "", "Template slug")
	meetingTemplatesCreateCmd.Flags().String("hosts", "", "Comma-separated host emails")
	meetingTemplatesCreateCmd.Flags().String("kind", "", "Meeting kind")
	meetingTemplatesCreateCmd.Flags().String("spot", "", "Meeting spot")
	meetingTemplatesCreateCmd.Flags().Int("duration", 0, "Meeting duration in minutes")
	meetingTemplatesCreateCmd.Flags().String("json-file", "", "Path to JSON file with template data")

	meetingTemplatesCmd.AddCommand(meetingTemplatesUpdateCmd)
	meetingTemplatesUpdateCmd.Flags().String("name", "", "Template name")
	meetingTemplatesUpdateCmd.Flags().String("slug", "", "Template slug")
	meetingTemplatesUpdateCmd.Flags().String("json-file", "", "Path to JSON file with template data")

	meetingTemplatesCmd.AddCommand(meetingTemplatesDeleteCmd)
}
