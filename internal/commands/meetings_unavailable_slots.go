package commands

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var meetingsUnavailableSlotsCmd = &cobra.Command{
	Use:   "unavailable-slots <meeting-sid>",
	Short: "Explain which times a meeting cannot be booked and why",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		params := meetingUnavailableSlotsParams(cmd)

		data, err := c.Get(fmt.Sprintf("/meetings/%s/unavailable-slots", args[0]), params)
		if err != nil {
			return err
		}

		printList(data, "dates", nil)
		return nil
	},
}

func meetingUnavailableSlotsParams(cmd *cobra.Command) url.Values {
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

	duration, _ := cmd.Flags().GetInt("duration")
	if duration > 0 {
		params.Set("duration", fmt.Sprintf("%d", duration))
	}

	return params
}

func init() {
	meetingsCmd.AddCommand(meetingsUnavailableSlotsCmd)

	meetingsUnavailableSlotsCmd.Flags().Int("year", 0, "Year")
	meetingsUnavailableSlotsCmd.Flags().Int("month", 0, "Month (1-12)")
	meetingsUnavailableSlotsCmd.Flags().Int("day", 0, "Day of month")
	meetingsUnavailableSlotsCmd.Flags().String("time-zone", "", timeZoneUsage)
	meetingsUnavailableSlotsCmd.Flags().String("host-email", "", "Only this host of a round-robin or multi-host meeting")
	meetingsUnavailableSlotsCmd.Flags().Int("duration", 0, "Slot length in minutes; must be one of the meeting's durations")

	markFlagsRequired(meetingsUnavailableSlotsCmd, "year", "month", "time-zone")
}
