package configuration

type PricesInfo struct {
	WeekDayHourlyPrice    float32
	HoursWeekDay          int
	WeekendDayHourlyPrice float32
	HoursWeekendDay       int
	BhDayHourlyPrice      float32
	HoursBhDay            int
}

func (c *Configuration) GetPricesInfo() (*PricesInfo, error) {
	weekDayPrice, err := c.FindPriceByDay("weekday")
	if err != nil {
		return nil, err
	}
	excludedWeekDayHoursAmount := 0
	excludedHours := c.FindRotationExcludedHoursByDay("weekday")
	if excludedHours != nil {
		excludedWeekDayHoursAmount = excludedHours.ExcludedEndsAt - excludedHours.ExcludedStartsAt
	}
	weekDayWorkingHours := 24 - excludedWeekDayHoursAmount

	weekendDayPrice, err := c.FindPriceByDay("weekend")
	if err != nil {
		return nil, err
	}
	excludedWeekendDayHoursAmount := 0
	excludedHours = c.FindRotationExcludedHoursByDay("weekend")
	if excludedHours != nil {
		excludedWeekendDayHoursAmount = excludedHours.ExcludedEndsAt - excludedHours.ExcludedStartsAt
	}
	weekendDayWorkingHours := 24 - excludedWeekendDayHoursAmount

	bhDayPrice, err := c.FindPriceByDay("bankholiday")
	if err != nil {
		return nil, err
	}
	excludedBhDayHoursAmount := 0
	excludedHours = c.FindRotationExcludedHoursByDay("bankholiday")
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
