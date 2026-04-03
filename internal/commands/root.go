package commands

import (
	"fmt"
	"os"

	"github.com/bigbinary/neeto-cal-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "neetocal",
	Short: "NeetoCal CLI — manage your calendar from the terminal",
	Long:  "A command-line interface for NeetoCal. Manage meetings, bookings, availabilities, and more.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		quietFlag, _ := cmd.Flags().GetBool("quiet")
		output.ForceJSON = jsonFlag
		output.QuietMode = quietFlag
	},
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "Output as JSON")
	rootCmd.PersistentFlags().Bool("quiet", false, "Output raw data only (no envelope)")
	rootCmd.PersistentFlags().String("subdomain", "", "Override saved subdomain")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
