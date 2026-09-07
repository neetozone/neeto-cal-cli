package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var automationRulesCmd = &cobra.Command{
	Use:   "automation-rules",
	Short: "Manage automation rules",
}

var automationRulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List automation rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		params := paginationParams(cmd)

		data, err := c.Get("/automation-rules", params)
		if err != nil {
			return err
		}

		printList(data, "automation_rules", nil)
		return nil
	},
}

var automationRulesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an automation rule",
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

		event, _ := cmd.Flags().GetString("event")
		if event != "" {
			body["event"] = event
		}

		data, err := c.Post("/automation-rules", body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var automationRulesDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an automation rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		err = c.Delete(fmt.Sprintf("/automation-rules/%s", args[0]))
		if err != nil {
			return err
		}

		printMessage("Automation rule deleted.")
		return nil
	},
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(automationRulesCmd) })

	automationRulesCmd.AddCommand(automationRulesListCmd)
	addPaginationFlags(automationRulesListCmd)

	automationRulesCmd.AddCommand(automationRulesCreateCmd)
	automationRulesCreateCmd.Flags().String("name", "", "Rule name")
	automationRulesCreateCmd.Flags().String("event", "", "Trigger event")
	automationRulesCreateCmd.Flags().String("json-file", "", "Path to JSON file with meeting_ids and actions")
	markFlagsRequired(automationRulesCreateCmd, "event", "json-file")
	allowJSONFileToSatisfyRequiredFlags(automationRulesCreateCmd)

	automationRulesCmd.AddCommand(automationRulesDeleteCmd)
}
