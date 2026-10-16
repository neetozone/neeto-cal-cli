package commands

import (
	"net/url"

	"github.com/spf13/cobra"
)

var holidaysCmd = &cobra.Command{
	Use:   "holidays",
	Short: "View workspace holidays",
}

var holidaysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspace holidays",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/holidays", holidaysListParams(cmd))
		if err != nil {
			return err
		}

		printList(data, "holidays", nil)
		return nil
	},
}

func holidaysListParams(cmd *cobra.Command) url.Values {
	params := paginationParams(cmd)

	for flag, param := range map[string]string{
		"search":    "search",
		"from-date": "from_date",
		"to-date":   "to_date",
	} {
		value, _ := cmd.Flags().GetString(flag)
		if value != "" {
			params.Set(param, value)
		}
	}

	if includeDisabled, _ := cmd.Flags().GetBool("include-disabled"); includeDisabled {
		params.Set("include_disabled", "true")
	}

	return params
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(holidaysCmd) })

	holidaysCmd.AddCommand(holidaysListCmd)
	addPaginationFlags(holidaysListCmd)
	holidaysListCmd.Flags().String("search", "", "Only holidays whose name contains this text")
	holidaysListCmd.Flags().String("from-date", "", "Only holidays on or after this date (YYYY-MM-DD); recurring holidays are always included")
	holidaysListCmd.Flags().String("to-date", "", "Only holidays on or before this date (YYYY-MM-DD)")
	holidaysListCmd.Flags().Bool("include-disabled", false, "Include holidays that are turned off")
}
