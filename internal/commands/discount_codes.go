package commands

import (
	"github.com/spf13/cobra"
)

var discountCodesCmd = &cobra.Command{
	Use:   "discount-codes",
	Short: "Manage discount codes",
}

var discountCodesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a discount code",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		code, _ := cmd.Flags().GetString("code")
		kind, _ := cmd.Flags().GetString("kind")
		value, _ := cmd.Flags().GetFloat64("value")

		body := map[string]interface{}{
			"code":  code,
			"kind":  kind,
			"value": value,
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

func init() {
	register(func(root *cobra.Command) { root.AddCommand(discountCodesCmd) })

	discountCodesCmd.AddCommand(discountCodesCreateCmd)
	discountCodesCreateCmd.Flags().String("code", "", "Discount code")
	discountCodesCreateCmd.Flags().String("kind", "", "Code kind (percentage, fixed)")
	discountCodesCreateCmd.Flags().Float64("value", 0, "Discount value: 0-100 when --kind is percentage, otherwise a flat amount in the currency of the meetings it applies to")
	discountCodesCreateCmd.Flags().String("meeting-ids", "", "Comma-separated meeting SIDs")
	discountCodesCreateCmd.Flags().String("expires-at", "", "Expiration date (YYYY-MM-DD)")
	markFlagsRequired(discountCodesCreateCmd, "code", "kind", "value")
}
