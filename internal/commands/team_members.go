package commands

import (
	"fmt"
	"net/url"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var teamMembersCmd = &cobra.Command{
	Use:   "team-members",
	Short: "Manage team members",
}

var teamMembersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List team members",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		params := teamMembersListParams(cmd)

		data, err := c.Get("/team-members", params)
		if err != nil {
			return err
		}

		breadcrumbs := []output.Breadcrumb{
			{Label: "Show details", Command: "neetocal team-members show <id>"},
		}

		printList(data, "team_members", breadcrumbs)
		return nil
	},
}

func teamMembersListParams(cmd *cobra.Command) url.Values {
	params := paginationParams(cmd)

	email, _ := cmd.Flags().GetString("email")
	if email != "" {
		params.Set("email", email)
	}

	return params
}

var teamMembersShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a team member",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/team-members/%s", args[0]), nil)
		if err != nil {
			return err
		}

		breadcrumbs := []output.Breadcrumb{
			{Label: "List team members", Command: "neetocal team-members list"},
			{Label: "Delete team member", Command: fmt.Sprintf("neetocal team-members delete %s", args[0])},
		}

		printResource(data, breadcrumbs)
		return nil
	},
}

var teamMembersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Add team members",
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

		emails, _ := cmd.Flags().GetString("emails")
		if emails != "" {
			body["emails"] = splitCSV(emails)
		}

		organizationRole, _ := cmd.Flags().GetString("organization-role")
		if organizationRole != "" {
			body["organization_role"] = organizationRole
		}

		invitedBy, _ := cmd.Flags().GetString("invited-by")
		if invitedBy != "" {
			body["invited_by"] = invitedBy
		}

		if cmd.Flags().Changed("send-invitation-email") {
			sendInvitationEmail, _ := cmd.Flags().GetBool("send-invitation-email")
			body["send_invitation_email"] = sendInvitationEmail
		}

		data, err := c.Post("/team-members", body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var teamMembersUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a team member",
	Args:  cobra.ExactArgs(1),
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

		email, _ := cmd.Flags().GetString("email")
		if email != "" {
			body["email"] = email
		}

		firstName, _ := cmd.Flags().GetString("first-name")
		if firstName != "" {
			body["first_name"] = firstName
		}

		lastName, _ := cmd.Flags().GetString("last-name")
		if lastName != "" {
			body["last_name"] = lastName
		}

		timeZone, _ := cmd.Flags().GetString("time-zone")
		if timeZone != "" {
			body["time_zone"] = timeZone
		}

		organizationRole, _ := cmd.Flags().GetString("organization-role")
		if organizationRole != "" {
			body["organization_role"] = organizationRole
		}

		data, err := c.Patch(fmt.Sprintf("/team-members/%s", args[0]), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var teamMembersDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a team member",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		err = c.Delete(fmt.Sprintf("/team-members/%s", args[0]))
		if err != nil {
			return err
		}

		printMessage("Team member removed.")
		return nil
	},
}

var teamMembersSlotsCmd = &cobra.Command{
	Use:   "slots",
	Short: "List start times when every listed person is free for the whole duration",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/common-slots", teamMembersSlotsParams(cmd))
		if err != nil {
			return err
		}

		printList(data, "slots", nil)
		return nil
	},
}

func teamMembersSlotsParams(cmd *cobra.Command) url.Values {
	emails, _ := cmd.Flags().GetString("emails")
	duration, _ := cmd.Flags().GetInt("duration")
	startDate, _ := cmd.Flags().GetString("start-date")
	endDate, _ := cmd.Flags().GetString("end-date")
	timeZone, _ := cmd.Flags().GetString("time-zone")

	params := url.Values{}
	for _, email := range splitCSV(emails) {
		params.Add("emails[]", email)
	}
	params.Set("duration", fmt.Sprintf("%d", duration))
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)
	params.Set("time_zone", timeZone)

	return params
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(teamMembersCmd) })

	teamMembersCmd.AddCommand(teamMembersListCmd)
	addPaginationFlags(teamMembersListCmd)
	teamMembersListCmd.Flags().String("email", "", "Filter by exact email")

	teamMembersCmd.AddCommand(teamMembersShowCmd)

	teamMembersCmd.AddCommand(teamMembersCreateCmd)
	teamMembersCreateCmd.Flags().String("emails", "", "Comma-separated emails to add")
	teamMembersCreateCmd.Flags().String("organization-role", "", "Organization role (defaults to Standard)")
	teamMembersCreateCmd.Flags().String("invited-by", "", "Email of an existing workspace admin")
	teamMembersCreateCmd.Flags().Bool("send-invitation-email", true, "Send an invitation email to the added members")
	teamMembersCreateCmd.Flags().String("json-file", "", "Path to JSON file with team member data")
	markFlagsRequired(teamMembersCreateCmd, "emails")
	allowJSONFileToSatisfyRequiredFlags(teamMembersCreateCmd)

	teamMembersCmd.AddCommand(teamMembersUpdateCmd)
	teamMembersUpdateCmd.Flags().String("email", "", "Email")
	teamMembersUpdateCmd.Flags().String("first-name", "", "First name")
	teamMembersUpdateCmd.Flags().String("last-name", "", "Last name")
	teamMembersUpdateCmd.Flags().String("time-zone", "", timeZoneUsage)
	teamMembersUpdateCmd.Flags().String("organization-role", "", "Organization role")
	teamMembersUpdateCmd.Flags().String("json-file", "", "Path to JSON file with team member data")

	teamMembersCmd.AddCommand(teamMembersDeleteCmd)

	teamMembersCmd.AddCommand(teamMembersSlotsCmd)
	teamMembersSlotsCmd.Flags().String("emails", "", "Comma-separated emails of the people who must all be free")
	teamMembersSlotsCmd.Flags().Int("duration", 0, "Meeting length in minutes")
	teamMembersSlotsCmd.Flags().String("start-date", "", "First date to search (YYYY-MM-DD)")
	teamMembersSlotsCmd.Flags().String("end-date", "", "Last date to search (YYYY-MM-DD, at most 31 days after --start-date)")
	teamMembersSlotsCmd.Flags().String("time-zone", "", timeZoneUsage)
	markFlagsRequired(teamMembersSlotsCmd, "emails", "duration", "start-date", "end-date", "time-zone")
}
