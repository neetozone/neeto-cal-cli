package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var packagesPurchasesCmd = &cobra.Command{
	Use:   "purchases",
	Short: "Manage package purchases",
}

var packagesPurchasesListCmd = &cobra.Command{
	Use:   "list <package-id>",
	Short: "List purchases for a package",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		params := paginationParams(cmd)

		data, err := c.Get(fmt.Sprintf("/packages/%s/purchases", args[0]), params)
		if err != nil {
			return err
		}

		printList(data, "purchases", nil)
		return nil
	},
}

var packagesPurchasesShowCmd = &cobra.Command{
	Use:   "show <package-id> <purchase-id>",
	Short: "Show a package purchase",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient()
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/packages/%s/purchases/%s", args[0], args[1]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

func init() {
	packagesCmd.AddCommand(packagesPurchasesCmd)

	packagesPurchasesCmd.AddCommand(packagesPurchasesListCmd)
	addPaginationFlags(packagesPurchasesListCmd)

	packagesPurchasesCmd.AddCommand(packagesPurchasesShowCmd)
}
