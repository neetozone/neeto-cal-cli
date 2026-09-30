package commands

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var meetingsSlotsCmd = &cobra.Command{
	Use:   "slots <meeting-sid>",
	Short: "List available slots for a meeting",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		params := meetingSlotsParams(cmd)

		data, err := c.Get(fmt.Sprintf("/meetings/%s/slots", args[0]), params)
		if err != nil {
			return err
		}

		printList(data, "slots", nil)
		return nil
	},
}

func meetingSlotsParams(cmd *cobra.Command) url.Values {
	year, _ := cmd.Flags().GetInt("year")
	month, _ := cmd.Flags().GetInt("month")
	timeZone, _ := cmd.Flags().GetString("time-zone")

	params := url.Values{}
	params.Set("year", fmt.Sprintf("%d", year))
	params.Set("month", fmt.Sprintf("%d", month))
	params.Set("time_zone", timeZone)

	day, _ := cmd.Flags().GetInt("day")
	if day > 0 {
		params.Set("day", fmt.Sprintf("%d", day))
	}

	hostEmail, _ := cmd.Flags().GetString("host-email")
	if hostEmail != "" {
		params.Set("host_email", hostEmail)
	}

	return params
}

func init() {
	meetingsCmd.AddCommand(meetingsSlotsCmd)

	meetingsSlotsCmd.Flags().Int("year", 0, "Year")
	meetingsSlotsCmd.Flags().Int("month", 0, "Month (1-12)")
	meetingsSlotsCmd.Flags().Int("day", 0, "Day of month")
	meetingsSlotsCmd.Flags().String("time-zone", "", timeZoneUsage)
	meetingsSlotsCmd.Flags().String("host-email", "", "Only slots for this host of a round-robin or multi-host meeting")

	markFlagsRequired(meetingsSlotsCmd, "year", "month", "time-zone")
}
