package commands

import (
	"fmt"

	"github.com/neetozone/neeto-cal-cli/internal/output"
	"github.com/spf13/cobra"
)

var meetingsCmd = &cobra.Command{
	Use:   "meetings",
	Short: "Manage meetings",
}

var meetingsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all meetings",
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

		search, _ := cmd.Flags().GetString("search")
		if search != "" {
			params.Set("search", search)
		}

		data, err := c.Get("/meetings", params)
		if err != nil {
			return err
		}

		breadcrumbs := []output.Breadcrumb{
			{Label: "Show details", Command: "neetocal meetings show <sid>"},
		}

		printList(data, "meetings", breadcrumbs)
		return nil
	},
}

var meetingsShowCmd = &cobra.Command{
	Use:   "show <sid>",
	Short: "Show a meeting",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/meetings/%s", args[0]), nil)
		if err != nil {
			return err
		}

		breadcrumbs := []output.Breadcrumb{
			{Label: "List all meetings", Command: "neetocal meetings list"},
			{Label: "Delete meeting", Command: fmt.Sprintf("neetocal meetings delete %s", args[0])},
		}

		printResource(data, breadcrumbs)
		return nil
	},
}

var meetingsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a meeting",
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

		slug, _ := cmd.Flags().GetString("slug")
		if slug != "" {
			body["slug"] = slug
		}

		hosts, _ := cmd.Flags().GetString("hosts")
		if hosts != "" {
			body["hosts"] = splitCSV(hosts)
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

		description, _ := cmd.Flags().GetString("description")
		if description != "" {
			body["description"] = description
		}

		data, err := c.Post("/meetings", body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var meetingsUpdateCmd = &cobra.Command{
	Use:   "update <sid>",
	Short: "Update a meeting",
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

		slug, _ := cmd.Flags().GetString("slug")
		if slug != "" {
			body["slug"] = slug
		}

		description, _ := cmd.Flags().GetString("description")
		if description != "" {
			body["description"] = description
		}

		hosts, _ := cmd.Flags().GetString("hosts")
		if hosts != "" {
			body["hosts"] = splitCSV(hosts)
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

		data, err := c.Put(fmt.Sprintf("/meetings/%s", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var meetingsDeleteCmd = &cobra.Command{
	Use:   "delete <sid>",
	Short: "Delete a meeting",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		err = c.Delete(fmt.Sprintf("/meetings/%s", args[0]))
		if err != nil {
			return err
		}

		output.PrintMessage("Meeting deleted successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(meetingsCmd)

	meetingsCmd.AddCommand(meetingsListCmd)
	addPaginationFlags(meetingsListCmd)
	meetingsListCmd.Flags().String("host-email", "", "Filter by host email")
	meetingsListCmd.Flags().String("search", "", "Search meetings")

	meetingsCmd.AddCommand(meetingsShowCmd)

	meetingsCmd.AddCommand(meetingsCreateCmd)
	meetingsCreateCmd.Flags().String("name", "", "Meeting name")
	meetingsCreateCmd.Flags().String("slug", "", "Meeting slug")
	meetingsCreateCmd.Flags().String("hosts", "", "Comma-separated host emails")
	meetingsCreateCmd.Flags().String("kind", "", "Meeting kind")
	meetingsCreateCmd.Flags().String("spot", "", "Meeting spot")
	meetingsCreateCmd.Flags().Int("duration", 0, "Meeting duration in minutes")
	meetingsCreateCmd.Flags().String("description", "", "Meeting description")
	meetingsCreateCmd.Flags().String("json-file", "", "Path to JSON file with meeting data")

	meetingsCmd.AddCommand(meetingsUpdateCmd)
	meetingsUpdateCmd.Flags().String("name", "", "Meeting name")
	meetingsUpdateCmd.Flags().String("slug", "", "Meeting slug")
	meetingsUpdateCmd.Flags().String("description", "", "Meeting description")
	meetingsUpdateCmd.Flags().String("hosts", "", "Comma-separated host emails (optional; omit to keep current hosts)")
	meetingsUpdateCmd.Flags().String("kind", "", "Meeting kind")
	meetingsUpdateCmd.Flags().String("spot", "", "Meeting spot")
	meetingsUpdateCmd.Flags().Int("duration", 0, "Meeting duration in minutes")
	meetingsUpdateCmd.Flags().String("json-file", "", "Path to JSON file with meeting data")

	meetingsCmd.AddCommand(meetingsDeleteCmd)
}
