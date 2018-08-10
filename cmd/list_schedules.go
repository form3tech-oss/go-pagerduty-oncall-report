package cmd

import (
	"fmt"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"
	"github.com/spf13/cobra"
)

var listSchedulesCmd = &cobra.Command{
	Use:   "schedules",
	Short: "list schedules on PagerDuty",
	Long:  "Get the list of schedules configured in PagerDuty",
	RunE:  listSchedules,
}

func init() {
	rootCmd.AddCommand(listSchedulesCmd)
}

func listSchedules(cmd *cobra.Command, args []string) error {
	schedules, err := api.Client.ListSchedules()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("==== Found %d schedule(s) ====", len(schedules)))
	for _, schedule := range schedules {
		fmt.Println(fmt.Sprintf("[%s] %-20s, Timezone: %s", schedule.ID, schedule.Name, schedule.TimeZone))
	}

	return nil
}
