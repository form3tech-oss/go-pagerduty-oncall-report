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
	log.Printf("startDate: %s, endDate: %s", startDate, endDate)

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

	weekDayHourlyPrice, weekendDayHourlyPrice, bhDayHourlyPrice, err := getPrices()
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Hourly prices - Week day: %f, Weekend day: %f, Bank holiday: %f",
		weekDayHourlyPrice, weekendDayHourlyPrice, bhDayHourlyPrice))

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

		scheduleData, err := generateScheduleData(scheduleInfo, usersRotationData,
			weekDayHourlyPrice, weekendDayHourlyPrice, bhDayHourlyPrice)
		if err != nil {
			return err
		}

		printableData.SchedulesData = append(printableData.SchedulesData, scheduleData)
	}

	summaryPrintableData := calculateSummaryData(printableData.SchedulesData)
	printableData.UsersSchedulesSummary = summaryPrintableData

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

func calculateSummaryData(data []*report.ScheduleData) []*report.UserSchedulesSummary {

	usersSummary := make(map[string]*report.UserSchedulesSummary, 0)

	for _, schedData := range data {
		for _, schedUser := range schedData.RotaUsers {
			userSummary, ok := usersSummary[schedUser.Name]
			if !ok {
				userSummary = &report.UserSchedulesSummary{
					Name: schedUser.Name,
				}
				usersSummary[schedUser.Name] = userSummary
			}

			userSummary.NumWorkHours += schedUser.NumWorkHours
			userSummary.NumWeekendHours += schedUser.NumWeekendHours
			userSummary.NumBankHolidaysHours += schedUser.NumBankHolidaysHours
			userSummary.TotalAmountWorkHours += schedUser.TotalAmountWorkHours
			userSummary.TotalAmountWeekendHours += schedUser.TotalAmountWeekendHours
			userSummary.TotalAmountBankHolidaysHours += schedUser.TotalAmountBankHolidaysHours
			userSummary.TotalAmount += schedUser.TotalAmount
		}
	}

	var result []*report.UserSchedulesSummary
	for _, userSummary := range usersSummary {
		result = append(result, userSummary)
	}

	return result
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

func generateScheduleData(scheduleInfo *api.ScheduleInfo, usersRotationData api.ScheduleUserRotationData,
	weekDayHourlyPrice, weekendDayHourlyPrice, bhDayHourlyPrice float32) (*report.ScheduleData, error) {

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
			currentDate := period.Start
			periodNumBHH, periodNumWendH, periodNumWH := float32(0), float32(0), float32(0)
			for currentDate.Before(period.End) {
				//log.Println(userRotaInfo.Name, "current date:", currentDate, "- bank holiday: ", userCalendar.IsDateBankHoliday(currentDate), "- weekend: ", userCalendar.IsWeekend(currentDate))
				if userCalendar.IsDateBankHoliday(currentDate) {
					periodNumBHH += 0.5
				} else if userCalendar.IsWeekend(currentDate) {
					periodNumWendH += 0.5
				} else {
					periodNumWH += 0.5
				}

				currentDate = currentDate.Add(time.Minute * 30)
			}
			scheduleUserData.NumBankHolidaysHours += periodNumBHH
			scheduleUserData.NumWeekendHours += periodNumWendH
			scheduleUserData.NumWorkHours += periodNumWH
			//fmt.Println(fmt.Sprintf("User rota info [%s] - start: %s, end: %s, work hours: %f, weekend hours: %f, bank holiday hours: %f",
			//	userRotaInfo.Name, period.Start, period.End, periodNumWH, periodNumWendH, periodNumBHH))
			periodNumBHH, periodNumWendH, periodNumWH = float32(0), float32(0), float32(0)
		}

		scheduleUserData.TotalAmountWorkHours = scheduleUserData.NumWorkHours * weekDayHourlyPrice
		scheduleUserData.TotalAmountWeekendHours = scheduleUserData.NumWeekendHours * weekendDayHourlyPrice
		scheduleUserData.TotalAmountBankHolidaysHours = scheduleUserData.NumBankHolidaysHours * bhDayHourlyPrice
		scheduleUserData.TotalAmount = scheduleUserData.TotalAmountWorkHours +
			scheduleUserData.TotalAmountWeekendHours +
			scheduleUserData.TotalAmountBankHolidaysHours
		scheduleData.RotaUsers = append(scheduleData.RotaUsers, scheduleUserData)
	}

	return scheduleData, nil
}

func getPrices() (float32, float32, float32, error) {

	weekDayPrice, err := Config.FindPriceByDay("weekday")
	if err != nil {
		return 0, 0, 0, err
	}
	weekendDayPrice, err := Config.FindPriceByDay("weekend")
	if err != nil {
		return 0, 0, 0, err
	}
	bhDayPrice, err := Config.FindPriceByDay("bankholiday")
	if err != nil {
		return 0, 0, 0, err
	}

	return float32(*weekDayPrice) / 24, float32(*weekendDayPrice) / 24, float32(*bhDayPrice) / 24, nil
}
