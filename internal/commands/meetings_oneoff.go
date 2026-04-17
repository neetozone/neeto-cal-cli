package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var meetingsOneOffLinkCmd = &cobra.Command{
	Use:   "one-off-link <meeting-sid>",
	Short: "Create a one-off link for a meeting",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Post(fmt.Sprintf("/meetings/%s/one-off-links", args[0]), nil)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func init() {
	meetingsCmd.AddCommand(meetingsOneOffLinkCmd)
}
