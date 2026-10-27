package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newBlockedEmailsCreateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create"}
	cmd.Flags().String("emails", "", "")
	cmd.Flags().String("note", "", "")
	return cmd
}

func TestBlockedEmailsCreateBody_SplitsEmailsAndSendsNote(t *testing.T) {
	cmd := newBlockedEmailsCreateTestCmd()
	setFlags(t, cmd, map[string]string{"emails": "one@example.com, two@example.com", "note": "Spam"})

	body := blockedEmailsCreateBody(cmd)

	want := map[string]interface{}{"blocked_email": map[string]interface{}{
		"emails": []string{"one@example.com", "two@example.com"},
		"note":   "Spam",
	}}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestBlockedEmailsCreateBody_OmitsUnsetNote(t *testing.T) {
	cmd := newBlockedEmailsCreateTestCmd()
	setFlags(t, cmd, map[string]string{"emails": "one@example.com"})

	body := blockedEmailsCreateBody(cmd)

	want := map[string]interface{}{"blocked_email": map[string]interface{}{"emails": []string{"one@example.com"}}}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}
