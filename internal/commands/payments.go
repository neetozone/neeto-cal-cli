package commands

import (
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var paymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "List payments for bookings and package purchases",
}

var paymentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List payments, newest first",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/payments", paymentFilterParams(cmd))
		if err != nil {
			return err
		}

		printList(data, "payments", nil)
		return nil
	},
}

var refundsCmd = &cobra.Command{
	Use:   "refunds",
	Short: "List refunds of booking and package purchase payments",
}

var refundsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List refunds, newest first",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/refunds", paymentFilterParams(cmd))
		if err != nil {
			return err
		}

		printList(data, "refunds", nil)
		return nil
	},
}

func paymentFilterParams(cmd *cobra.Command) url.Values {
	params := paginationParams(cmd)
	for _, flag := range []string{"from-time", "to-time", "status", "type", "meeting-sid"} {
		if value, _ := cmd.Flags().GetString(flag); value != "" {
			params.Set(strings.ReplaceAll(flag, "-", "_"), value)
		}
	}
	return params
}

func addPaymentFilterFlags(cmd *cobra.Command, records, statusHelp string) {
	cmd.Flags().String("from-time", "", "Only "+records+" made at or after this time (ISO 8601 datetime or YYYY-MM-DD)")
	cmd.Flags().String("to-time", "", "Only "+records+" made before this time (ISO 8601 datetime or YYYY-MM-DD)")
	cmd.Flags().String("status", "", statusHelp)
	cmd.Flags().String("type", "", "Only payments for this: booking or package_purchase")
	cmd.Flags().String("meeting-sid", "", "Only bookings of this meeting, by SID or slug")
	addPaginationFlags(cmd)
}

func init() {
	register(func(root *cobra.Command) {
		root.AddCommand(paymentsCmd)
		root.AddCommand(refundsCmd)
	})

	paymentsCmd.AddCommand(paymentsListCmd)
	addPaymentFilterFlags(paymentsListCmd, "payments", "Payment status (default: successful), e.g. declined, pending")

	refundsCmd.AddCommand(refundsListCmd)
	addPaymentFilterFlags(refundsListCmd, "refunds", "Refund status (default: successful), e.g. pending, declined")
}
