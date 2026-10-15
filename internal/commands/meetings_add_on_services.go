package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var meetingsAddOnServicesCmd = &cobra.Command{
	Use:   "add-on-services",
	Short: "Manage meeting add-on services",
}

var meetingsAddOnServicesListCmd = &cobra.Command{
	Use:   "list <meeting-sid>",
	Short: "List add-on services for a meeting",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/meetings/%s/add-on-services", args[0]), nil)
		if err != nil {
			return err
		}

		printList(data, "add_on_services", nil)
		return nil
	},
}

var meetingsAddOnServicesShowCmd = &cobra.Command{
	Use:   "show <meeting-sid> <id>",
	Short: "Show a meeting add-on service",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/meetings/%s/add-on-services/%s", args[0], args[1]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var meetingsAddOnServicesCreateCmd = &cobra.Command{
	Use:   "create <meeting-sid>",
	Short: "Create a meeting add-on service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Post(fmt.Sprintf("/meetings/%s/add-on-services", args[0]), addOnServiceBody(cmd))
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var meetingsAddOnServicesUpdateCmd = &cobra.Command{
	Use:   "update <meeting-sid> <id>",
	Short: "Update a meeting add-on service",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Put(fmt.Sprintf("/meetings/%s/add-on-services/%s", args[0], args[1]), addOnServiceBody(cmd))
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func addOnServiceBody(cmd *cobra.Command) map[string]interface{} {
	body := map[string]interface{}{}

	for _, flag := range []string{"name", "description", "price"} {
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetString(flag)
			body[flag] = value
		}
	}

	for _, flag := range []string{"duration", "position"} {
		if cmd.Flags().Changed(flag) {
			value, _ := cmd.Flags().GetInt(flag)
			body[flag] = value
		}
	}

	return body
}

var meetingsAddOnServicesDeleteCmd = &cobra.Command{
	Use:   "delete <meeting-sid> <id>",
	Short: "Delete a meeting add-on service",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		err = c.Delete(fmt.Sprintf("/meetings/%s/add-on-services/%s", args[0], args[1]))
		if err != nil {
			return err
		}

		printMessage("Add-on service deleted.")
		return nil
	},
}

func addAddOnServiceFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Add-on service name")
	cmd.Flags().String("description", "", "Add-on service description")
	cmd.Flags().Int("duration", 0, "Extra minutes added to the booking")
	cmd.Flags().Int("position", 0, "Position in the list, starting at 1")
	cmd.Flags().String("price", "", "Price such as 25.00; 0 makes it free")
}

func init() {
	meetingsCmd.AddCommand(meetingsAddOnServicesCmd)

	meetingsAddOnServicesCmd.AddCommand(meetingsAddOnServicesListCmd)
	meetingsAddOnServicesCmd.AddCommand(meetingsAddOnServicesShowCmd)

	meetingsAddOnServicesCmd.AddCommand(meetingsAddOnServicesCreateCmd)
	addAddOnServiceFlags(meetingsAddOnServicesCreateCmd)
	markFlagsRequired(meetingsAddOnServicesCreateCmd, "name")

	meetingsAddOnServicesCmd.AddCommand(meetingsAddOnServicesUpdateCmd)
	addAddOnServiceFlags(meetingsAddOnServicesUpdateCmd)

	meetingsAddOnServicesCmd.AddCommand(meetingsAddOnServicesDeleteCmd)
}
