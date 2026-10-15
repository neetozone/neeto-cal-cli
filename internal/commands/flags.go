package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

const timeZoneUsage = "IANA time zone name, e.g. America/New_York. Full list: " +
	"https://en.wikipedia.org/wiki/List_of_tz_database_time_zones"

func setJSONObjectFlag(cmd *cobra.Command, body map[string]interface{}, flag, key string) error {
	raw, _ := cmd.Flags().GetString(flag)
	if raw == "" {
		return nil
	}

	var object map[string]any
	if err := json.Unmarshal([]byte(raw), &object); err != nil {
		return fmt.Errorf("--%s must be a JSON object: %w", flag, err)
	}

	body[key] = object
	return nil
}
