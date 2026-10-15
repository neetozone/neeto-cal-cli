package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newBookingsPaymentsCreateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create"}
	for _, flag := range []string{"payment-provider", "identifier", "discount-code", "idempotency-key", "amount", "currency"} {
		cmd.Flags().String(flag, "", "")
	}
	cmd.Flags().Float64("tip-value", 0, "")
	cmd.Flags().Bool("tip-is-custom", false, "")
	return cmd
}

func newBookingsPaymentsUpdateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "update"}
	for _, flag := range []string{"payment-provider", "status", "notes", "vpa-address", "performer-id"} {
		cmd.Flags().String(flag, "", "")
	}
	return cmd
}

func TestBookingPaymentCreateBody_SendsCashAmountAndIdempotencyKey(t *testing.T) {
	cmd := newBookingsPaymentsCreateTestCmd()
	setFlags(t, cmd, map[string]string{
		"payment-provider": "cash",
		"idempotency-key":  "key-1",
		"amount":           "150.00",
		"currency":         "USD",
	})

	body := bookingPaymentCreateBody(cmd)

	want := map[string]interface{}{
		"payment_provider": "cash",
		"idempotency_key":  "key-1",
		"amount":           "150.00",
		"currency":         "USD",
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestBookingPaymentCreateBody_SendsTipAttributesOnlyWhenTipValueIsSet(t *testing.T) {
	cmd := newBookingsPaymentsCreateTestCmd()
	setFlags(t, cmd, map[string]string{"payment-provider": "stripe", "tip-value": "10", "tip-is-custom": "true"})

	body := bookingPaymentCreateBody(cmd)

	want := map[string]interface{}{
		"payment_provider": "stripe",
		"tip_attributes":   map[string]interface{}{"value": 10.0, "is_custom": true},
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}

	cmd = newBookingsPaymentsCreateTestCmd()
	setFlags(t, cmd, map[string]string{"payment-provider": "stripe", "tip-is-custom": "true"})

	if _, present := bookingPaymentCreateBody(cmd)["tip_attributes"]; present {
		t.Errorf("tip_attributes present without --tip-value")
	}
}

func TestBookingPaymentUpdateBody_SendsVpaAddressAndPerformerID(t *testing.T) {
	cmd := newBookingsPaymentsUpdateTestCmd()
	setFlags(t, cmd, map[string]string{
		"payment-provider": "upi",
		"status":           "successful",
		"vpa-address":      "oliver@upi",
		"performer-id":     "member-1",
	})

	body := bookingPaymentUpdateBody(cmd)

	want := map[string]interface{}{
		"payment_provider": "upi",
		"status":           "successful",
		"vpa_address":      "oliver@upi",
		"performer_id":     "member-1",
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}
