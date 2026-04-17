package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bookingsPaymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "Manage booking payments",
}

var bookingsPaymentsCreateCmd = &cobra.Command{
	Use:   "create <booking-id>",
	Short: "Create a payment for a booking",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		paymentProvider, _ := cmd.Flags().GetString("payment-provider")

		body := map[string]interface{}{
			"payment_provider": paymentProvider,
		}

		identifier, _ := cmd.Flags().GetString("identifier")
		if identifier != "" {
			body["identifier"] = identifier
		}

		discountCode, _ := cmd.Flags().GetString("discount-code")
		if discountCode != "" {
			body["discount_code"] = discountCode
		}

		data, err := c.Post(fmt.Sprintf("/bookings/%s/payments", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var bookingsPaymentsUpdateCmd = &cobra.Command{
	Use:   "update <booking-id> <payment-id>",
	Short: "Update a booking payment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		paymentProvider, _ := cmd.Flags().GetString("payment-provider")
		status, _ := cmd.Flags().GetString("status")

		body := map[string]interface{}{
			"payment_provider": paymentProvider,
			"status":           status,
		}

		notes, _ := cmd.Flags().GetString("notes")
		if notes != "" {
			body["notes"] = notes
		}

		data, err := c.Put(fmt.Sprintf("/bookings/%s/payments/%s", args[0], args[1]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func init() {
	bookingsCmd.AddCommand(bookingsPaymentsCmd)

	bookingsPaymentsCmd.AddCommand(bookingsPaymentsCreateCmd)
	bookingsPaymentsCreateCmd.Flags().String("payment-provider", "", "Payment provider")
	bookingsPaymentsCreateCmd.Flags().String("identifier", "", "Payment identifier")
	bookingsPaymentsCreateCmd.Flags().String("discount-code", "", "Discount code")
	_ = bookingsPaymentsCreateCmd.MarkFlagRequired("payment-provider")

	bookingsPaymentsCmd.AddCommand(bookingsPaymentsUpdateCmd)
	bookingsPaymentsUpdateCmd.Flags().String("payment-provider", "", "Payment provider")
	bookingsPaymentsUpdateCmd.Flags().String("status", "", "Payment status (successful, rejected)")
	bookingsPaymentsUpdateCmd.Flags().String("notes", "", "Payment notes")
	_ = bookingsPaymentsUpdateCmd.MarkFlagRequired("payment-provider")
	_ = bookingsPaymentsUpdateCmd.MarkFlagRequired("status")
}
