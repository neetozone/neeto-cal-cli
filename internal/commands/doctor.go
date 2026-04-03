package commands

import (
	"fmt"
	"net/http"
	"time"

	"github.com/neetozone/neeto-cal-cli/internal/auth"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check CLI health and connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check credentials
		creds, err := auth.LoadCredentials()
		if err != nil {
			fmt.Println("✗ Authentication: not logged in")
			fmt.Println("  Run 'neetocal login' to authenticate")
			return nil
		}
		fmt.Printf("✓ Authentication: logged in as %s\n", creds.Email)

		// Check connectivity
		baseURL := auth.BaseURL(creds.Subdomain)
		httpClient := &http.Client{Timeout: 10 * time.Second}
		start := time.Now()
		resp, err := httpClient.Get(baseURL)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("✗ API connection: could not reach %s\n", baseURL)
			fmt.Printf("  Error: %v\n", err)
		} else {
			resp.Body.Close()
			fmt.Printf("✓ API connection: %s (responding in %dms)\n", baseURL, elapsed.Milliseconds())
		}

		// Check version
		fmt.Printf("✓ CLI version: %s\n", Version)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
