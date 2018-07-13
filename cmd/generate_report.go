package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/rogersole/go-pagerduty-oncall-report/api"
	"github.com/rogersole/go-pagerduty-oncall-report/configuration"
	"github.com/rogersole/go-pagerduty-oncall-report/report"
	"github.com/spf13/cobra"
)

var (
	scheduleReportCmd = &cobra.Command{
		Use:   "report",
		Short: "generates the report(s) for the given schedule(s) id(s)",
		Long:  "Generates the report of the given list of schedules or all (except the ignored ones configured in yml)",
		RunE:  generateReport,
	}

	schedules    []string
	outputFormat string
)

func init() {
	scheduleReportCmd.Flags().StringSliceVarP(&schedules, "schedules", "s", []string{"all"}, "schedule ids to report (comma-separated with no spaces), or 'all'")
	scheduleReportCmd.Flags().StringVarP(&outputFormat, "output-format", "o", "console", "pdf, console")
	rootCmd.AddCommand(scheduleReportCmd)
}

type InputData struct {
	schedules []string
	startDate time.Time
	endDate   time.Time
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func processArguments() InputData {

	if !contains([]string{"console", "pdf"}, outputFormat) {
		log.Printf("output format %s not supported. Defaulting to 'console'", outputFormat)
		outputFormat = "console"
	}

	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	startDate := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	if len(schedules) == 1 && schedules[0] == "all" {
		schedules = []string{}
		schedulesList, err := api.Client.ListSchedules()
		if err != nil {
			log.Fatalln(fmt.Sprintf("Error getting the schedules list: %s", err.Error()))
		}
		for _, schedule := range schedulesList {
			if !Config.IsScheduleIDToIgnore(schedule.ID) {
				schedules = append(schedules, schedule.ID)
			} else {
				log.Println(fmt.Sprintf("Ignoring schedule '%s'", schedule.ID))
			}
		}
	}

	return InputData{
		schedules: schedules,
		startDate: startDate,
		endDate:   endDate,
	}
}

func generateReport(cmd *cobra.Command, args []string) error {
	input := processArguments()

	configuration.LoadCalendars(input.startDate.Year())
	printableData := &report.PrintableData{
		Start:         input.startDate,
		End:           input.endDate,
		SchedulesData: make([]*report.ScheduleData, 0),
	}

	for _, scheduleID := range input.schedules {
		log.Printf("Loading information for the schedule '%s'", scheduleID)
		scheduleInfo, err := getScheduleInformation(scheduleID, input.startDate, input.endDate)
		if err != nil {
			return err
		}

		usersRotationData, err := getUsersRotationData(scheduleInfo)
		if err != nil {
			return err
		}

		scheduleData, err := generateScheduleData(scheduleInfo, usersRotationData)
		if err != nil {
			return err
		}

		printableData.SchedulesData = append(printableData.SchedulesData, scheduleData)
	}

	var reportWriter report.Writer
	if outputFormat == "pdf" {
		reportWriter = report.NewPDFReport(Config.Currency)
	} else {
		reportWriter = report.NewConsoleReport(Config.Currency)
	}
	message, err := reportWriter.GenerateReport(printableData)
	if err != nil {
		return err
	}

	if len(message) > 0 {
		log.Println(message)
	}
	return nil
}

func getScheduleInformation(scheduleID string, startDate, endDate time.Time) (*api.ScheduleInfo, error) {
	schedule, err := api.Client.GetSchedule(scheduleID,
		startDate.Format("2006-01-02T15:04:05"),
		endDate.Format("2006-01-02T15:04:05"))
	if err != nil {
		return nil, err
	}

	location, _ := time.LoadLocation(schedule.TimeZone)

	scheduleInfo := &api.ScheduleInfo{
		ID:            scheduleID,
		Name:          schedule.Name,
		Location:      location,
		Start:         startDate,
		End:           endDate,
		FinalSchedule: schedule.FinalSchedule,
	}
	return scheduleInfo, nil
}

func getUsersRotationData(scheduleInfo *api.ScheduleInfo) (api.ScheduleUserRotationData, error) {
	usersInfo := api.ScheduleUserRotationData{}
	for _, entry := range scheduleInfo.FinalSchedule.RenderedScheduleEntries {
		startDate, err := time.ParseInLocation(time.RFC3339, entry.Start, scheduleInfo.Location)
		if err != nil {
			return nil, err
		}
		endDate, err := time.ParseInLocation(time.RFC3339, entry.End, scheduleInfo.Location)
		if err != nil {
			return nil, err
		}
		//fmt.Println(fmt.Sprintf("[%s] %-25s %v - %v", entry.User.ID, entry.User.Summary, startDate, endDate))

		userRotaInfo, ok := usersInfo[entry.User.ID]
		if !ok {
			userRotaInfo = &api.UserRotaInfo{
				ID:      entry.User.ID,
				Name:    entry.User.Summary,
				Periods: make([]*api.UserRotaPeriod, 0),
			}
			usersInfo[entry.User.ID] = userRotaInfo
		}
		newUserRotaPeriod := &api.UserRotaPeriod{
			Start: startDate,
			End:   endDate,
		}

		userRotaInfo.Periods = append(userRotaInfo.Periods, newUserRotaPeriod)
	}

	return usersInfo, nil
}

func generateScheduleData(scheduleInfo *api.ScheduleInfo, usersRotationData api.ScheduleUserRotationData) (*report.ScheduleData, error) {
	weekDayPrice, err := Config.FindPriceByDay("weekday")
	if err != nil {
		return nil, err
	}
	weekendDayPrice, err := Config.FindPriceByDay("weekend")
	if err != nil {
		return nil, err
	}
	bhDayPrice, err := Config.FindPriceByDay("bankholiday")
	if err != nil {
		return nil, err
	}

	scheduleData := &report.ScheduleData{
		ID:        scheduleInfo.ID,
		Name:      scheduleInfo.Name,
		RotaUsers: make([]*report.ScheduleUser, 0),
	}
	for userID, userRotaInfo := range usersRotationData {
		rotationUserConfig, err := Config.FindRotationUserInfoByID(userID)
		if err != nil {
			return nil, err
		}

		calendarName := fmt.Sprintf("%s-%d", rotationUserConfig.HolidaysCalendar, scheduleInfo.Start.Year())
		userCalendar, present := configuration.BankHolidaysCalendars[calendarName]
		if !present {
			return nil, fmt.Errorf("calendar '%s' not found for user '%s'. Aborting", calendarName, userID)
		}

		scheduleUserData := &report.ScheduleUser{
			Name: userRotaInfo.Name,
		}

		for _, period := range userRotaInfo.Periods {
			//fmt.Println(fmt.Sprintf("User rota info [%s] - start: %s, end: %s", userID, period.Start, period.End))
			currentDate := period.Start
			for currentDate.Before(period.End) {
				//log.Println(userRotaInfo.Name, "current date:", currentDate, "- bank holiday: ", userCalendar.IsDateBankHoliday(currentDate), "- weekend: ", userCalendar.IsWeekend(currentDate))
				if userCalendar.IsDateBankHoliday(currentDate) {
					scheduleUserData.NumBankHolidaysDays++
				} else if userCalendar.IsWeekend(currentDate) {
					scheduleUserData.NumWeekendDays++
				} else {
					scheduleUserData.NumWorkDays++
				}
				currentDate = currentDate.Add(time.Hour * 24)
			}
		}

		scheduleUserData.TotalAmountWorkDays = scheduleUserData.NumWorkDays * weekDayPrice
		scheduleUserData.TotalAmountWeekendDays = scheduleUserData.NumWeekendDays * weekendDayPrice
		scheduleUserData.TotalAmountBankHolidaysDays = scheduleUserData.NumBankHolidaysDays * bhDayPrice
		scheduleUserData.TotalAmount = scheduleUserData.TotalAmountWorkDays +
			scheduleUserData.TotalAmountWeekendDays +
			scheduleUserData.TotalAmountBankHolidaysDays
		scheduleData.RotaUsers = append(scheduleData.RotaUsers, scheduleUserData)
	}

	return scheduleData, nil
}
