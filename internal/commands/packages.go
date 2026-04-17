package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "Manage packages",
}

var packagesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		params := paginationParams(cmd)

		data, err := c.Get("/packages", params)
		if err != nil {
			return err
		}

		printList(data, "packages", nil)
		return nil
	},
}

var packagesShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a package",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/packages/%s", args[0]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(packagesCmd)

	packagesCmd.AddCommand(packagesListCmd)
	addPaginationFlags(packagesListCmd)

	packagesCmd.AddCommand(packagesShowCmd)
}
