package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/neetozone/neeto-cal-cli/internal/auth"
	"github.com/neetozone/neeto-cal-cli/internal/output"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to NeetoCal via browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		subdomain, _ := cmd.Flags().GetString("subdomain")

		if subdomain == "" {
			// Check if already saved
			creds, err := auth.LoadCredentials()
			if err == nil && creds.Subdomain != "" {
				subdomain = creds.Subdomain
			} else {
				fmt.Print("Enter your NeetoCal subdomain (e.g., 'acme' for acme.neetocal.com): ")
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				subdomain = strings.TrimSpace(input)
			}
		}

		if subdomain == "" {
			return fmt.Errorf("subdomain is required")
		}

		creds, err := auth.Login(subdomain)
		if err != nil {
			return err
		}

		output.PrintMessage(fmt.Sprintf("Logged in as %s", creds.Email))
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and clear saved credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.ClearCredentials(); err != nil {
			return err
		}
		output.PrintMessage("Logged out successfully.")
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current authenticated user",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := auth.LoadCredentials()
		if err != nil {
			return err
		}

		output.PrintMessage(fmt.Sprintf("Logged in as %s on %s.neetocal.com", creds.Email, creds.Subdomain))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
}
