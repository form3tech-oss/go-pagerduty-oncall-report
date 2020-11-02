package configuration

import (
	"fmt"
	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"
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

func (c *Configuration) FindRotationExcludedHoursByDay(dayType string) (*RotationExcludedHoursDay, error) {
	if excludedInfo, ok := c.cacheExcludedByDay[dayType]; ok {
		return excludedInfo, nil
	}

	for _, rotationExcludedHours := range c.RotationExcludedHours {
		if rotationExcludedHours.Day == dayType {
			c.cacheExcludedByDay[rotationExcludedHours.Day] = &rotationExcludedHours
			return &rotationExcludedHours, nil
		}
	}

	return nil, fmt.Errorf("day type %s not found", dayType)
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

	user, err := api.Client.GetUserById(userID)
	if err != nil {
		return nil, fmt.Errorf("could not find user from pagerduty api, err: %v", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user id %s not found", userID)
	}

	rotationUser := &RotationUser{
		UserID:           userID,
		Name:             user.Name,
		HolidaysCalendar: c.DefaultHolidayCalendar, // default to uk
	}

	c.cacheRotationUsers[userID] = rotationUser
	log.Printf("defaulting user %s id: %s to %s\n", user.Name, userID, c.DefaultHolidayCalendar)

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
