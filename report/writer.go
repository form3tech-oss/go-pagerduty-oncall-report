package report

import (
	"time"
)

type PrintableData struct {
	Start                 time.Time
	End                   time.Time
	SchedulesData         []*ScheduleData
	UsersSchedulesSummary []*ScheduleUser
}

type ScheduleData struct {
	ID        string
	Name      string
	StartDate time.Time
	EndDate   time.Time
	RotaUsers []*ScheduleUser
}

type ScheduleUser struct {
	Name                         string
	EmailAddress                 string
	NumWorkHours                 float32
	NumWorkDays                  float32
	TotalAmountWorkHours         float32
	NumWeekendHours              float32
	NumWeekendDays               float32
	TotalAmountWeekendHours      float32
	NumBankHolidaysHours         float32
	NumBankHolidaysDays          float32
	TotalAmountBankHolidaysHours float32
	TotalAmount                  float32
}

type Writer interface {
	GenerateReport(data *PrintableData) (string, error)
}
