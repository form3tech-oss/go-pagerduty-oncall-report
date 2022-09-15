package configuration

import (
	"fmt"
	"log"
)

type RotationUser struct {
	UserID           string
	Name             string
	HolidaysCalendar string
}

type RotationPriceDay struct {
	Day   string
	Price int
}

type RotationPrices struct {
	Currency string
	DaysInfo []RotationPriceDay
}

type RotationExcludedHoursDay struct {
	Day              string
	ExcludedStartsAt int
	ExcludedEndsAt   int
}

type RotationInfo struct {
	DailyRotationStartsAt    int
	CheckRotationChangeEvery int
}

type ReportTimeRange struct {
	Start string
	End   string
}

type ScheduleTimeRange struct {
	Id    string
	Start string
	End   string
}

type Configuration struct {
	PdAuthToken                string
	DefaultHolidayCalendar     string
	DefaultUserTimezone        string
	ReportTimeRange            ReportTimeRange
	RotationInfo               RotationInfo
	RotationExcludedHours      []RotationExcludedHoursDay
	RotationPrices             RotationPrices
	RotationUsers              []RotationUser
	ScheduleTimeRangeOverrides []ScheduleTimeRange
	SchedulesToIgnore          []string

	cacheRotationUsers  map[string]*RotationUser
	cacheRotationPrices map[string]int
	cacheExcludedByDay  map[string]*RotationExcludedHoursDay
}

func New() *Configuration {
	return &Configuration{
		cacheRotationUsers:  make(map[string]*RotationUser),
		cacheRotationPrices: make(map[string]int),
		cacheExcludedByDay:  make(map[string]*RotationExcludedHoursDay),
	}
}

func (c *Configuration) FindPriceByDay(dayType string) (*int, error) {
	if price, ok := c.cacheRotationPrices[dayType]; ok {
		return &price, nil
	}

	for _, rotationPrice := range c.RotationPrices.DaysInfo {
		if rotationPrice.Day == dayType {
			c.cacheRotationPrices[rotationPrice.Day] = rotationPrice.Price
			return &rotationPrice.Price, nil
		}
	}

	return nil, fmt.Errorf("day type %s not found", dayType)
}

func (c *Configuration) FindRotationExcludedHoursByDay(dayType string) *RotationExcludedHoursDay {
	if excludedInfo, ok := c.cacheExcludedByDay[dayType]; ok {
		return excludedInfo
	}

	for _, rotationExcludedHours := range c.RotationExcludedHours {
		if rotationExcludedHours.Day == dayType {
			c.cacheExcludedByDay[rotationExcludedHours.Day] = &rotationExcludedHours
			return &rotationExcludedHours
		}
	}

	return nil
}

func (c *Configuration) FindRotationUserInfoByID(userID string) (*RotationUser, error) {
	if rotationUser, ok := c.cacheRotationUsers[userID]; ok {
		return rotationUser, nil
	}

	for _, rotationUser := range c.RotationUsers {
		if rotationUser.UserID == userID {
			c.cacheRotationUsers[userID] = &rotationUser
			return &rotationUser, nil
		}
	}
	if c.DefaultHolidayCalendar == "" { // if you dont specify a default calendar fallback to old behaviour
		return nil, fmt.Errorf("user id %s not found", userID)
	}

	rotationUser := &RotationUser{
		UserID:           userID,
		HolidaysCalendar: c.DefaultHolidayCalendar, // default to config value
	}

	c.cacheRotationUsers[userID] = rotationUser
	log.Printf("defaulting user with id: %s to %s\n", userID, c.DefaultHolidayCalendar)

	return rotationUser, nil
}

func (c *Configuration) IsScheduleIDToIgnore(scheduleID string) bool {
	for _, scheduleIDToIgnore := range c.SchedulesToIgnore {
		if scheduleIDToIgnore == scheduleID {
			return true
		}
	}
	return false
}
