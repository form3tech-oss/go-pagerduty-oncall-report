package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"log"
	"time"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"
	"github.com/form3tech-oss/go-pagerduty-oncall-report/configuration"
	"github.com/form3tech-oss/go-pagerduty-oncall-report/report"
	"github.com/spf13/cobra"
)

var (
	scheduleReportCmd = &cobra.Command{
		Use:   "report",
		Short: "generates the report(s) for the given schedule(s) id(s)",
		Long:  "Generates the report of the given list of schedules or all (except the ignored ones configured in yml)",
		RunE:  generateReport,
	}

	rawSchedules []string
	outputFormat string
	directory    string
)

func init() {
	scheduleReportCmd.Flags().StringSliceVarP(&rawSchedules, "schedules", "s", []string{"all"}, "schedule ids to report (comma-separated with no spaces), or 'all'")
	scheduleReportCmd.Flags().StringVarP(&outputFormat, "output-format", "o", "console", "pdf, console, csv")
	scheduleReportCmd.Flags().StringVarP(&directory, "output", "d", "", "output path (default is $HOME)")
	rootCmd.AddCommand(scheduleReportCmd)
}

type Schedule struct {
	id        string
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

func processArguments() []Schedule {

	if !contains([]string{"console", "pdf", "csv"}, outputFormat) {
		log.Printf("output format %s not supported. Defaulting to 'console'", outputFormat)
		outputFormat = "console"
	}
	if directory == "" {
		directory, _ = homedir.Dir()
	}
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)

	var defaultStartDate time.Time
	if Config.ReportTimeRange.Start != "" {
		var err error
		defaultStartDate, err = time.Parse(time.RFC822, Config.ReportTimeRange.Start)
		if err != nil {
			log.Fatalf("Error parsing report start time: %s", err)
		}
	} else {
		defaultStartDate = time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	}

	var defaultEndDate time.Time
	if Config.ReportTimeRange.End != "" {
		var err error
		defaultEndDate, err = time.Parse(time.RFC822, Config.ReportTimeRange.End)
		if err != nil {
			log.Fatalf("Error parsing report end time: %s", err)
		}
	} else {
		defaultEndDate = defaultStartDate.AddDate(0, 1, 0)
		defaultEndDate = defaultEndDate.Add(time.Hour * time.Duration(Config.RotationInfo.DailyRotationStartsAt))
	}

	startOverrides := make(map[string]time.Time)
	endOverrides := make(map[string]time.Time)

	for _, override := range Config.ScheduleTimeRangeOverrides {
		var err error
		startOverrides[override.Id], err = time.Parse(time.RFC822, override.Start)
		if err != nil {
			log.Fatalf("Error parsing start time override for schedule %s: %v", override.Id, err)
		}
		endOverrides[override.Id], err = time.Parse(time.RFC822, override.End)
		if err != nil {
			log.Fatalf("Error parsing end time override for schedule %s: %v", override.Id, err)
		}
	}

	schedules := make([]Schedule, 0)
	if len(rawSchedules) == 1 && rawSchedules[0] == "all" {
		schedulesList, err := api.Client.ListSchedules()
		if err != nil {
			log.Fatalln(fmt.Sprintf("Error getting the schedules list: %s", err.Error()))
		}

		for _, schedule := range schedulesList {
			if !Config.IsScheduleIDToIgnore(schedule.ID) {
				var thisStartDate time.Time
				if _, ok := startOverrides[schedule.ID]; ok {
					thisStartDate = startOverrides[schedule.ID]
				} else {
					thisStartDate = defaultStartDate
				}

				var thisEndDate time.Time
				if _, ok := endOverrides[schedule.ID]; ok {
					thisEndDate = endOverrides[schedule.ID]
				} else {
					thisEndDate = defaultEndDate
				}

				schedules = append(schedules, Schedule{
					id:        schedule.ID,
					startDate: thisStartDate,
					endDate:   thisEndDate,
				})

				log.Printf("[%s] defaultStartDate: %s, defaultEndDate: %s", schedule.ID, thisStartDate, thisEndDate)

			} else {
				log.Println(fmt.Sprintf("Ignoring schedule '%s'", schedule.ID))
			}
		}
	} else {
		for _, schedule := range rawSchedules {
			if !Config.IsScheduleIDToIgnore(schedule) {
				var thisStartDate time.Time
				if _, ok := startOverrides[schedule]; ok {
					thisStartDate = startOverrides[schedule]
				} else {
					thisStartDate = defaultStartDate
				}

				var thisEndDate time.Time
				if _, ok := endOverrides[schedule]; ok {
					thisEndDate = endOverrides[schedule]
				} else {
					thisEndDate = defaultEndDate
				}

				schedules = append(schedules, Schedule{
					id:        schedule,
					startDate: thisStartDate,
					endDate:   thisEndDate,
				})

				log.Printf("[%s] defaultStartDate: %s, defaultEndDate: %s", schedule, thisStartDate, thisEndDate)
			} else {
				log.Fatal("Configuration explicitly ignores schedule passed as parameter - check your config.")
			}
		}
	}

	return schedules
}

func generateReport(cmd *cobra.Command, args []string) error {
	input := processArguments()

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
	printableData := &report.PrintableData{
		Start:         firstStartDate,
		End:           lastEndDate,
		SchedulesData: make([]*report.ScheduleData, 0),
	}

	pricesInfo, err := generatePricesInfo()
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Hourly prices (in %s) - Week day: %v (%vh), Weekend day: %v (%vh), Bank holiday: %v (%vh)",
		Config.RotationPrices.Currency, pricesInfo.WeekDayHourlyPrice, pricesInfo.HoursWeekDay, pricesInfo.WeekendDayHourlyPrice,
		pricesInfo.HoursWeekendDay, pricesInfo.BhDayHourlyPrice, pricesInfo.HoursBhDay))

	for _, schedule := range input {
		log.Printf("Loading information for the schedule '%s'", schedule.id)
		scheduleInfo, err := getScheduleInformation(schedule.id, schedule.startDate, schedule.endDate)
		if err != nil {
			return err
		}

		usersRotationData, err := getUsersRotationData(scheduleInfo)
		if err != nil {
			return err
		}

		scheduleData, err := generateScheduleData(scheduleInfo, usersRotationData, pricesInfo, schedule)
		if err != nil {
			return err
		}

		printableData.SchedulesData = append(printableData.SchedulesData, scheduleData)
	}

	summaryPrintableData := calculateSummaryData(printableData.SchedulesData, pricesInfo)
	printableData.UsersSchedulesSummary = summaryPrintableData

	var reportWriter report.Writer
	if outputFormat == "pdf" {
		reportWriter = report.NewPDFReport(Config.RotationPrices.Currency, directory)
	} else if outputFormat == "csv" {
		reportWriter = report.NewCsvReport(Config.RotationPrices.Currency, directory)
	} else {
		reportWriter = report.NewConsoleReport(Config.RotationPrices.Currency)
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

func calculateSummaryData(data []*report.ScheduleData, pricesInfo *PricesInfo) []*report.ScheduleUser {

	usersSummary := make(map[string]*report.ScheduleUser)

	for _, schedData := range data {
		for _, schedUser := range schedData.RotaUsers {
			userSummary, ok := usersSummary[schedUser.Name]
			if !ok {
				userSummary = &report.ScheduleUser{
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

	result := make([]*report.ScheduleUser, 0)
	for _, userSummary := range usersSummary {
		userSummary.NumWorkDays = userSummary.NumWorkHours / float32(pricesInfo.HoursWeekDay)
		userSummary.NumWeekendDays = userSummary.NumWeekendHours / float32(pricesInfo.HoursWeekendDay)
		userSummary.NumBankHolidaysDays = userSummary.NumBankHolidaysHours / float32(pricesInfo.HoursBhDay)
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
	pricesInfo *PricesInfo, schedule Schedule) (*report.ScheduleData, error) {

	scheduleData := &report.ScheduleData{
		ID:        scheduleInfo.ID,
		Name:      scheduleInfo.Name,
		StartDate: schedule.startDate,
		EndDate:   schedule.endDate,
		RotaUsers: make([]*report.ScheduleUser, 0),
	}

	for userID, userRotaInfo := range usersRotationData {
		rotationUserConfig, err := Config.FindRotationUserInfoByID(userID)
		if err != nil {
			log.Println("Error:", err)
			continue
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
			currentMonth := period.Start.Month()
			currentDate := period.Start
			for currentDate.Before(period.End) {
				updateDataForDate(&userCalendar, scheduleUserData, currentMonth, currentDate)
				currentDate = currentDate.Add(time.Minute * time.Duration(Config.RotationInfo.CheckRotationChangeEvery))
			}
		}

		scheduleUserData.NumWorkDays = scheduleUserData.NumWorkHours / float32(pricesInfo.HoursWeekDay)
		scheduleUserData.NumWeekendDays = scheduleUserData.NumWeekendHours / float32(pricesInfo.HoursWeekendDay)
		scheduleUserData.NumBankHolidaysDays = scheduleUserData.NumBankHolidaysHours / float32(pricesInfo.HoursBhDay)
		scheduleUserData.TotalAmountWorkHours = scheduleUserData.NumWorkHours * pricesInfo.WeekDayHourlyPrice
		scheduleUserData.TotalAmountWeekendHours = scheduleUserData.NumWeekendHours * pricesInfo.WeekendDayHourlyPrice
		scheduleUserData.TotalAmountBankHolidaysHours = scheduleUserData.NumBankHolidaysHours * pricesInfo.BhDayHourlyPrice
		scheduleUserData.TotalAmount = scheduleUserData.TotalAmountWorkHours +
			scheduleUserData.TotalAmountWeekendHours +
			scheduleUserData.TotalAmountBankHolidaysHours
		scheduleData.RotaUsers = append(scheduleData.RotaUsers, scheduleUserData)
	}

	return scheduleData, nil
}

func updateDataForDate(calendar *configuration.BHCalendar, data *report.ScheduleUser, currentMonth time.Month, date time.Time) {

	if date.Hour() < Config.RotationInfo.DailyRotationStartsAt {
		newDate := date.Add(time.Hour * time.Duration(-(date.Hour() + 1))) // move to yesterday night to determine which kind of day it was
		// if yesterday night was last month, ignore the date
		if newDate.Month() == currentMonth {
			updateDataForDate(calendar, data, currentMonth, newDate)
		}
	} else {
		if calendar.IsDateBankHoliday(date) {
			excludedHours, _ := Config.FindRotationExcludedHoursByDay("bankholiday")
			if excludedHours == nil {
				//fmt.Printf("%s - Month: %d, time: %v -- bank holiday\n", data.Name, currentMonth, date)
				data.NumBankHolidaysHours += 0.5
				return
			}

			if date.Hour() < excludedHours.ExcludedEndsAt && date.Hour() >= excludedHours.ExcludedEndsAt {
				//fmt.Printf("%s - Month: %d, time: %v -- bank holiday non excluded hours\n", data.Name, currentMonth, date)
				data.NumBankHolidaysHours += 0.5
			}
		} else if calendar.IsWeekend(date) {
			excludedHours, _ := Config.FindRotationExcludedHoursByDay("weekend")
			if excludedHours == nil {
				//fmt.Printf("%s - Month: %d, time: %v -- weekend\n", data.Name, currentMonth, date)
				data.NumWeekendHours += 0.5
				return
			}

			if date.Hour() < excludedHours.ExcludedEndsAt && date.Hour() >= excludedHours.ExcludedEndsAt {
				//fmt.Printf("%s - Month: %d, time: %v -- weekend non excluded hours\n", data.Name, currentMonth, date)
				data.NumWeekendHours += 0.5
			}
		} else {
			excludedHours, _ := Config.FindRotationExcludedHoursByDay("weekday")
			if excludedHours == nil {
				//fmt.Printf("%s - Month: %d, time: %v -- weekday\n", data.Name, currentMonth, date)
				data.NumWorkHours += 0.5
				return
			}

			if date.Hour() < excludedHours.ExcludedStartsAt || date.Hour() >= excludedHours.ExcludedEndsAt {
				//fmt.Printf("%s - Month: %d, time: %v -- weekday non excluded hours\n", data.Name, currentMonth, date)
				data.NumWorkHours += 0.5
			}
		}
	}
}

type PricesInfo struct {
	WeekDayHourlyPrice    float32
	HoursWeekDay          int
	WeekendDayHourlyPrice float32
	HoursWeekendDay       int
	BhDayHourlyPrice      float32
	HoursBhDay            int
}

func generatePricesInfo() (*PricesInfo, error) {

	weekDayPrice, err := Config.FindPriceByDay("weekday")
	if err != nil {
		return nil, err
	}
	excludedWeekDayHoursAmount := 0
	excludedHours, _ := Config.FindRotationExcludedHoursByDay("weekday")
	if excludedHours != nil {
		excludedWeekDayHoursAmount = excludedHours.ExcludedEndsAt - excludedHours.ExcludedStartsAt
	}
	weekDayWorkingHours := 24 - excludedWeekDayHoursAmount

	weekendDayPrice, err := Config.FindPriceByDay("weekend")
	if err != nil {
		return nil, err
	}
	excludedWeekendDayHoursAmount := 0
	excludedHours, _ = Config.FindRotationExcludedHoursByDay("weekend")
	if excludedHours != nil {
		excludedWeekendDayHoursAmount = excludedHours.ExcludedEndsAt - excludedHours.ExcludedStartsAt
	}
	weekendDayWorkingHours := 24 - excludedWeekendDayHoursAmount

	bhDayPrice, err := Config.FindPriceByDay("bankholiday")
	if err != nil {
		return nil, err
	}
	excludedBhDayHoursAmount := 0
	excludedHours, _ = Config.FindRotationExcludedHoursByDay("bankholiday")
	if excludedHours != nil {
		excludedBhDayHoursAmount = excludedHours.ExcludedEndsAt - excludedHours.ExcludedStartsAt
	}
	bhWorkingHours := 24 - excludedBhDayHoursAmount

	return &PricesInfo{
		WeekDayHourlyPrice:    float32(*weekDayPrice) / float32(weekDayWorkingHours),
		HoursWeekDay:          weekDayWorkingHours,
		WeekendDayHourlyPrice: float32(*weekendDayPrice) / float32(weekendDayWorkingHours),
		HoursWeekendDay:       weekendDayWorkingHours,
		BhDayHourlyPrice:      float32(*bhDayPrice) / float32(bhWorkingHours),
		HoursBhDay:            bhWorkingHours,
	}, nil
}
