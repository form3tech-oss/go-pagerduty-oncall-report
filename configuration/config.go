package configuration

import (
	"fmt"
)

type RotationUser struct {
	UserID           string
	Name             string
	HolidaysCalendar string
}

type RotationPrice struct {
	Type  string
	Price int
}

type Configuration struct {
	PdAuthToken       string
	RotationStartHour string
	RotationUsers     []RotationUser
	RotationPrices    []RotationPrice
	Currency          string
	SchedulesToIgnore []string

	cacheRotationUsers  map[string]*RotationUser
	cacheRotationPrices map[string]int
}

func New() *Configuration {
	return &Configuration{
		cacheRotationUsers:  make(map[string]*RotationUser),
		cacheRotationPrices: make(map[string]int),
	}
}

func (c *Configuration) FindPriceByDay(dayType string) (*int, error) {
	if price, ok := c.cacheRotationPrices[dayType]; ok {
		return &price, nil
	}

	for _, rotationPrice := range c.RotationPrices {
		if rotationPrice.Type == dayType {
			c.cacheRotationPrices[rotationPrice.Type] = rotationPrice.Price
			return &rotationPrice.Price, nil
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
