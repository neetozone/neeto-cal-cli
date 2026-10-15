package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var prefillTokensCmd = &cobra.Command{
	Use:   "prefill-tokens",
	Short: "Manage booking form prefill tokens",
}

var prefillTokensCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a prefill token for a booking form",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		body, err := prefillTokenBody(cmd)
		if err != nil {
			return err
		}

		data, err := c.Post("/prefill-tokens", body)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

func prefillTokenBody(cmd *cobra.Command) (map[string]interface{}, error) {
	jsonFile, _ := cmd.Flags().GetString("json-file")
	raw, _ := cmd.Flags().GetString("data")

	switch {
	case jsonFile != "":
		fileData, err := readJSONFile(jsonFile)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{"data": fileData}, nil
	case raw != "":
		var object map[string]any
		if err := json.Unmarshal([]byte(raw), &object); err != nil {
			return nil, fmt.Errorf("--data must be a JSON object: %w", err)
		}
		return map[string]interface{}{"data": object}, nil
	}

	return nil, errors.New("pass --data or --json-file")
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(prefillTokensCmd) })

	prefillTokensCmd.AddCommand(prefillTokensCreateCmd)
	prefillTokensCreateCmd.Flags().String("data", "", "Booking form answers as a JSON object keyed by question field code")
	prefillTokensCreateCmd.Flags().String("json-file", "", "Path to a JSON object used as the booking form answers")
}
