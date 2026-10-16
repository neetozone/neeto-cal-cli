package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newMeetingOutcomeTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("color", "", "")
	return cmd
}

func TestMeetingOutcomeBody_WrapsOnlyChangedFlags(t *testing.T) {
	cmd := newMeetingOutcomeTestCmd()
	setFlags(t, cmd, map[string]string{"color": "red"})

	body := meetingOutcomeBody(cmd)

	want := map[string]interface{}{"meeting_outcome": map[string]interface{}{"color": "red"}}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestMeetingOutcomeBody_SendsNameAndColor(t *testing.T) {
	cmd := newMeetingOutcomeTestCmd()
	setFlags(t, cmd, map[string]string{"name": "No show", "color": "gray"})

	body := meetingOutcomeBody(cmd)

	want := map[string]interface{}{"meeting_outcome": map[string]interface{}{"name": "No show", "color": "gray"}}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestMeetingOutcomesListParams_MapsSearchToSearchTerm(t *testing.T) {
	cmd := &cobra.Command{Use: "list"}
	cmd.Flags().String("search", "", "")
	setFlags(t, cmd, map[string]string{"search": "show"})

	if got := meetingOutcomesListParams(cmd).Get("search_term"); got != "show" {
		t.Errorf("search_term = %q, want show", got)
	}
}
