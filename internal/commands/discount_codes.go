package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var discountCodesCmd = &cobra.Command{
	Use:   "discount-codes",
	Short: "Manage discount codes",
}

var discountCodesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List discount codes",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/discount-codes", paginationParams(cmd))
		if err != nil {
			return err
		}

		printList(data, "discount_codes", nil)
		return nil
	},
}

var discountCodesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a discount code",
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

		code, _ := cmd.Flags().GetString("code")
		if code != "" {
			body["code"] = code
		}

		kind, _ := cmd.Flags().GetString("kind")
		if kind != "" {
			body["kind"] = kind
		}

		if cmd.Flags().Changed("value") {
			value, _ := cmd.Flags().GetInt("value")
			body["value"] = value
		}

		meetingIDs, _ := cmd.Flags().GetString("meeting-ids")
		if meetingIDs != "" {
			body["meeting_ids"] = splitCSV(meetingIDs)
		}

		expiresAt, _ := cmd.Flags().GetString("expires-at")
		if expiresAt != "" {
			body["expires_at"] = expiresAt
		}

		data, err := c.Post("/discount-codes", body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var discountCodesUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a discount code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		body, err := discountCodeUpdateBody(cmd)
		if err != nil {
			return err
		}

		data, err := c.Patch(fmt.Sprintf("/discount-codes/%s", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func discountCodeUpdateBody(cmd *cobra.Command) (map[string]interface{}, error) {
	body := map[string]interface{}{}

	jsonFile, _ := cmd.Flags().GetString("json-file")
	if jsonFile != "" {
		fileData, err := readJSONFile(jsonFile)
		if err != nil {
			return nil, err
		}
		body = fileData
	}

	for flag, field := range map[string]string{
		"code":       "code",
		"kind":       "kind",
		"expires-at": "expires_at",
	} {
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetString(flag)
			body[field] = value
		}
	}

	if cmd.Flags().Changed("value") {
		value, _ := cmd.Flags().GetInt("value")
		body["value"] = value
	}

	if cmd.Flags().Changed("meeting-ids") {
		meetingIDs, _ := cmd.Flags().GetString("meeting-ids")
		body["meeting_ids"] = splitCSV(meetingIDs)
	}

	return body, nil
}

var discountCodesDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a discount code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		if err := c.Delete(fmt.Sprintf("/discount-codes/%s", args[0])); err != nil {
			return err
		}

		printMessage("Discount code deleted.")
		return nil
	},
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(discountCodesCmd) })

	discountCodesCmd.AddCommand(discountCodesListCmd)
	addPaginationFlags(discountCodesListCmd)

	discountCodesCmd.AddCommand(discountCodesCreateCmd)
	discountCodesCreateCmd.Flags().String("code", "", "Discount code")
	discountCodesCreateCmd.Flags().String("kind", "", "Code kind (percentage, fixed)")
	discountCodesCreateCmd.Flags().Int("value", 0, "Discount value: a whole number, 1-100 when --kind is percentage, otherwise a flat amount in the workspace default currency")
	discountCodesCreateCmd.Flags().String("meeting-ids", "", "Comma-separated meeting SIDs")
	discountCodesCreateCmd.Flags().String("expires-at", "", "Expiration date (YYYY-MM-DD)")
	discountCodesCreateCmd.Flags().String("json-file", "", "Path to JSON file with discount code data")
	markFlagsRequired(discountCodesCreateCmd, "code", "kind", "value")
	allowJSONFileToSatisfyRequiredFlags(discountCodesCreateCmd)

	discountCodesCmd.AddCommand(discountCodesUpdateCmd)
	discountCodesUpdateCmd.Flags().String("code", "", "Discount code")
	discountCodesUpdateCmd.Flags().String("kind", "", "Code kind (percentage, fixed)")
	discountCodesUpdateCmd.Flags().Int("value", 0, "Discount value: a whole number, 1-100 when --kind is percentage, otherwise a flat amount in the workspace default currency")
	discountCodesUpdateCmd.Flags().String("meeting-ids", "", "Comma-separated meeting SIDs")
	discountCodesUpdateCmd.Flags().String("expires-at", "", "Expiration date (YYYY-MM-DD)")
	discountCodesUpdateCmd.Flags().String("json-file", "", "Path to JSON file with discount code data")

	discountCodesCmd.AddCommand(discountCodesDeleteCmd)
}
