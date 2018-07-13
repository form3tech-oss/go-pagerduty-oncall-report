package report

import (
	"fmt"
	"time"
)

type consoleReport struct {
	currency string
}

const (
	blankLine = ""
	separator = " -------------------------------------------------------------------------------------------------------------------------------"
	rowFormat = "| %-25s || %7v | %7v | %12v | %13v | %13v | %18v | %8v |"
)

func NewConsoleReport(currency string) Writer {
	return &consoleReport{
		currency: currency,
	}
}

func (r *consoleReport) GenerateReport(data *PrintableData) (string, error) {

	fmt.Println(separator)
	fmt.Println(fmt.Sprintf("| Generating report(s) from '%s' to '%s'", data.Start.Format("Mon Jan _2 15:04:05 2006"), data.End.Add(time.Second*-1).Format("Mon Jan _2 15:04:05 2006")))
	fmt.Println(separator)

	for _, scheduleData := range data.SchedulesData {
		fmt.Println(blankLine)
		fmt.Println(separator)
		fmt.Println(fmt.Sprintf("| Schedule: '%s' (%s)", scheduleData.Name, scheduleData.ID))
		fmt.Println(separator)
		fmt.Println(fmt.Sprintf(rowFormat, "USER", "WEEKDAY", "WEEKEND", "BANK HOLIDAY", "TOTAL WEEKDAY", "TOTAL WEEKEND", "TOTAL BANK HOLIDAY", "TOTAL"))
		fmt.Println(separator)

		for _, userData := range scheduleData.RotaUsers {
			fmt.Println(fmt.Sprintf(rowFormat, userData.Name,
				userData.NumWorkDays, userData.NumWeekendDays, userData.NumBankHolidaysDays,
				userData.TotalAmountWorkDays, userData.TotalAmountWeekendDays, userData.TotalAmountBankHolidaysDays,
				fmt.Sprintf("%s%d", r.currency, userData.TotalAmount)))
		}
	}

	fmt.Println(separator)

	return "", nil
}
