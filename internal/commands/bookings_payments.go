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

		body := bookingPaymentCreateBody(cmd)

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

		body := bookingPaymentUpdateBody(cmd)

		data, err := c.Put(fmt.Sprintf("/bookings/%s/payments/%s", args[0], args[1]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func bookingPaymentCreateBody(cmd *cobra.Command) map[string]interface{} {
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

	for flag, key := range map[string]string{
		"idempotency-key": "idempotency_key",
		"amount":          "amount",
		"currency":        "currency",
	} {
		value, _ := cmd.Flags().GetString(flag)
		if value != "" {
			body[key] = value
		}
	}

	if cmd.Flags().Changed("tip-value") {
		tipValue, _ := cmd.Flags().GetFloat64("tip-value")
		tipIsCustom, _ := cmd.Flags().GetBool("tip-is-custom")
		body["tip_attributes"] = map[string]interface{}{
			"value":     tipValue,
			"is_custom": tipIsCustom,
		}
	}

	return body
}

func bookingPaymentUpdateBody(cmd *cobra.Command) map[string]interface{} {
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

	for flag, key := range map[string]string{
		"vpa-address":  "vpa_address",
		"performer-id": "performer_id",
	} {
		value, _ := cmd.Flags().GetString(flag)
		if value != "" {
			body[key] = value
		}
	}

	return body
}

func init() {
	bookingsCmd.AddCommand(bookingsPaymentsCmd)

	bookingsPaymentsCmd.AddCommand(bookingsPaymentsCreateCmd)
	bookingsPaymentsCreateCmd.Flags().String("payment-provider", "", "Payment provider")
	bookingsPaymentsCreateCmd.Flags().String("identifier", "", "Payment identifier")
	bookingsPaymentsCreateCmd.Flags().String("discount-code", "", "Discount code")
	bookingsPaymentsCreateCmd.Flags().String("idempotency-key", "", "Key that makes a retried request return the same payment")
	bookingsPaymentsCreateCmd.Flags().String("amount", "", "Amount paid, such as 150.00 (cash only, for a meeting with no cash price; defaults to the meeting's price)")
	bookingsPaymentsCreateCmd.Flags().String("currency", "", "Three-letter currency code for --amount, such as USD")
	bookingsPaymentsCreateCmd.Flags().Float64("tip-value", 0, "Tip as a percentage of the amount")
	bookingsPaymentsCreateCmd.Flags().Bool("tip-is-custom", false, "Mark the tip as a custom amount chosen by the client")
	markFlagsRequired(bookingsPaymentsCreateCmd, "payment-provider")

	bookingsPaymentsCmd.AddCommand(bookingsPaymentsUpdateCmd)
	bookingsPaymentsUpdateCmd.Flags().String("payment-provider", "", "Payment provider")
	bookingsPaymentsUpdateCmd.Flags().String("status", "", "Payment status (successful, rejected)")
	bookingsPaymentsUpdateCmd.Flags().String("notes", "", "Payment notes")
	bookingsPaymentsUpdateCmd.Flags().String("vpa-address", "", "UPI VPA address the payment was received on (UPI payments marked successful)")
	bookingsPaymentsUpdateCmd.Flags().String("performer-id", "", "ID of the team member confirming the payment (ignored when logged in as a user)")
	markFlagsRequired(bookingsPaymentsUpdateCmd, "payment-provider", "status")
}
