package report

import "time"

type PrintableData struct {
	Start         time.Time
	End           time.Time
	SchedulesData []*ScheduleData
}

type ScheduleData struct {
	ID        string
	Name      string
	RotaUsers []*ScheduleUser
}

type ScheduleUser struct {
	Name                        string
	NumWorkDays                 int
	TotalAmountWorkDays         int
	NumWeekendDays              int
	TotalAmountWeekendDays      int
	NumBankHolidaysDays         int
	TotalAmountBankHolidaysDays int
	TotalAmount                 int
}

type Writer interface {
	GenerateReport(data *PrintableData) (string, error)
}
