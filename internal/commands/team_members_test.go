package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func newTeamMembersListTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "list"}
	addPaginationFlags(cmd)
	cmd.Flags().String("email", "", "Filter by exact email")
	return cmd
}

func TestTeamMembersListParams_SendsPageAndPageNumber(t *testing.T) {
	cmd := newTeamMembersListTestCmd()
	if err := cmd.Flags().Set("page", "3"); err != nil {
		t.Fatalf("set page: %v", err)
	}
	if err := cmd.Flags().Set("page-size", "50"); err != nil {
		t.Fatalf("set page-size: %v", err)
	}

	params := teamMembersListParams(cmd)

	if got := params.Get("page"); got != "3" {
		t.Errorf("page = %q, want 3", got)
	}
	if got := params.Get("page_number"); got != "3" {
		t.Errorf("page_number = %q, want 3", got)
	}
	if got := params.Get("page_size"); got != "50" {
		t.Errorf("page_size = %q, want 50", got)
	}
}

func TestTeamMembersListParams_Email(t *testing.T) {
	cmd := newTeamMembersListTestCmd()
	if err := cmd.Flags().Set("email", "oliver@example.com"); err != nil {
		t.Fatalf("set email: %v", err)
	}

	params := teamMembersListParams(cmd)

	if got := params.Get("email"); got != "oliver@example.com" {
		t.Errorf("email = %q, want oliver@example.com", got)
	}
}

func TestTeamMembersListParams_Empty(t *testing.T) {
	cmd := newTeamMembersListTestCmd()

	params := teamMembersListParams(cmd)

	if len(params) != 0 {
		t.Errorf("params = %v, want empty", params)
	}
}
