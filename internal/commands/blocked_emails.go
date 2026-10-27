package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var blockedEmailsCmd = &cobra.Command{
	Use:   "blocked-emails",
	Short: "Manage blocked emails",
}

var blockedEmailsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List blocked emails",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		params := paginationParams(cmd)
		search, _ := cmd.Flags().GetString("search")
		if search != "" {
			params.Set("search_term", search)
		}

		data, err := c.Get("/blocked-emails", params)
		if err != nil {
			return err
		}

		printList(data, "blocked_emails", nil)
		return nil
	},
}

var blockedEmailsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a blocked email",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/blocked-emails/%s", args[0]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var blockedEmailsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Block one or more emails",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Post("/blocked-emails", blockedEmailsCreateBody(cmd))
		if err != nil {
			return err
		}

		printList(data, "blocked_emails", nil)
		return nil
	},
}

func blockedEmailsCreateBody(cmd *cobra.Command) map[string]interface{} {
	emails, _ := cmd.Flags().GetString("emails")
	blockedEmail := map[string]interface{}{"emails": splitCSV(emails)}

	if cmd.Flags().Changed("note") {
		note, _ := cmd.Flags().GetString("note")
		blockedEmail["note"] = note
	}

	return map[string]interface{}{"blocked_email": blockedEmail}
}

var blockedEmailsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update the note on a blocked email",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		note, _ := cmd.Flags().GetString("note")
		body := map[string]interface{}{"blocked_email": map[string]interface{}{"note": note}}

		data, err := c.Put(fmt.Sprintf("/blocked-emails/%s", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var blockedEmailsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Unblock an email",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		err = c.Delete(fmt.Sprintf("/blocked-emails/%s", args[0]))
		if err != nil {
			return err
		}

		printMessage("Email unblocked.")
		return nil
	},
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(blockedEmailsCmd) })

	blockedEmailsCmd.AddCommand(blockedEmailsListCmd)
	blockedEmailsListCmd.Flags().String("search", "", "Search blocked emails by address")
	addPaginationFlags(blockedEmailsListCmd)

	blockedEmailsCmd.AddCommand(blockedEmailsShowCmd)

	blockedEmailsCmd.AddCommand(blockedEmailsCreateCmd)
	blockedEmailsCreateCmd.Flags().String("emails", "", "Comma-separated emails to block")
	blockedEmailsCreateCmd.Flags().String("note", "", "Why these emails are blocked")
	markFlagsRequired(blockedEmailsCreateCmd, "emails")

	blockedEmailsCmd.AddCommand(blockedEmailsUpdateCmd)
	blockedEmailsUpdateCmd.Flags().String("note", "", "New note; pass \"\" to clear it")
	markFlagsRequired(blockedEmailsUpdateCmd, "note")

	blockedEmailsCmd.AddCommand(blockedEmailsDeleteCmd)
}
