package report

import "time"

type PrintableData struct {
	Start                 time.Time
	End                   time.Time
	SchedulesData         []*ScheduleData
	UsersSchedulesSummary []*UserSchedulesSummary
}

type ScheduleData struct {
	ID        string
	Name      string
	RotaUsers []*ScheduleUser
}

type UserSchedulesSummary struct {
	Name                         string
	NumWorkHours                 float32
	TotalAmountWorkHours         float32
	NumWeekendHours              float32
	TotalAmountWeekendHours      float32
	NumBankHolidaysHours         float32
	TotalAmountBankHolidaysHours float32
	TotalAmount                  float32
}

type ScheduleUser struct {
	Name                         string
	NumWorkHours                 float32
	TotalAmountWorkHours         float32
	NumWeekendHours              float32
	TotalAmountWeekendHours      float32
	NumBankHolidaysHours         float32
	TotalAmountBankHolidaysHours float32
	TotalAmount                  float32
}

type Writer interface {
	GenerateReport(data *PrintableData) (string, error)
}
