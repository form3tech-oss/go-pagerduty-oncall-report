package cmd

import (
	"log"
	"time"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"
	"github.com/form3tech-oss/go-pagerduty-oncall-report/configuration"
	"github.com/form3tech-oss/go-pagerduty-oncall-report/report"

	"github.com/spf13/cobra"
)

var (
	schedule99dReportCmd = &cobra.Command{
		Use:   "report99d",
		Short: "generates the report(s) for the given schedule(s) id(s)",
		Long:  "Generates the report of the given list of schedules or all (except the ignored ones configured in yml)",
		RunE: func(cmd *cobra.Command, args []string) error {
			pd := &pagerDutyClient{
				client:              api.NewPagerDutyAPIClient(Config.PdAuthToken),
				defaultUserTimezone: Config.DefaultUserTimezone,
			}
			return pd.generateReport99d()
		},
	}
)

func init() {
	schedule99dReportCmd.Flags().StringSliceVarP(&rawSchedules, "schedules", "s", []string{"all"}, "schedule ids to report (comma-separated with no spaces), or 'all'")
	schedule99dReportCmd.Flags().StringVarP(&outputFormat, "output-format", "o", "console", "pdf, console, csv")
	schedule99dReportCmd.Flags().StringVarP(&directory, "output", "d", "", "output path (default is $HOME)")
	rootCmd.AddCommand(schedule99dReportCmd)
}

func (pd *pagerDutyClient) generateReport99d() error {
	input := pd.processArguments()
	firstStartDate := time.Now()
	lastEndDate := time.Time{}
	for _, schedule := range input {
		if schedule.startDate.Before(firstStartDate) {
			firstStartDate = schedule.startDate
		}
		if schedule.endDate.After(lastEndDate) {
			lastEndDate = schedule.endDate
		}
	}
	configuration.LoadCalendars(firstStartDate.Year())

	pricesInfo, err := Config.GetPricesInfo()
	if err != nil {
		return err
	}

	printableData := &report.PrintableData{
		Start:         firstStartDate,
		End:           lastEndDate,
		SchedulesData: make([]*report.ScheduleData, 0, len(input)),
	}
	for _, schedule := range input {
		log.Printf("Loading information for the schedule '%s'", schedule.id)
		scheduleInfo, err := pd.getScheduleInformation(schedule.id, schedule.startDate, schedule.endDate)
		if err != nil {
			return err
		}
		usersRotationData, err := getUsersRotationData(scheduleInfo)
		if err != nil {
			return err
		}
		scheduleData, err := pd.generateScheduleData(scheduleInfo, usersRotationData, pricesInfo, schedule)
		if err != nil {
			return err
		}
		printableData.SchedulesData = append(printableData.SchedulesData, scheduleData)
	}
	summaryPrintableData := calculateSummaryData(printableData.SchedulesData, pricesInfo)
	printableData.UsersSchedulesSummary = summaryPrintableData
	reportWriter := report.NewCsvOncallReport(directory)
	message, err := reportWriter.GenerateReport(printableData)
	if err != nil {
		return err
	}
	if len(message) > 0 {
		log.Println(message)
	}
	return nil
}
