package configuration

import (
	"fmt"
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

type Configuration struct {
	PdAuthToken           string
	RotationInfo          RotationInfo
	RotationExcludedHours []RotationExcludedHoursDay
	RotationPrices        RotationPrices
	RotationUsers         []RotationUser
	SchedulesToIgnore     []string

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
	return nil, fmt.Errorf("user id %s not found", userID)
}

func (c *Configuration) IsScheduleIDToIgnore(scheduleID string) bool {
	for _, scheduleIDToIgnore := range c.SchedulesToIgnore {
		if scheduleIDToIgnore == scheduleID {
			return true
		}
	}
	return false
}
