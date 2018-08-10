package stages

import (
	"testing"

	"bytes"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/configuration"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type ConfigStage struct {
	t *testing.T

	configRaw            []byte
	config               *configuration.Configuration
	configError          error
	configUnmarshalError error

	mapValue interface{}
	mapError error
}

func ConfigTest(t *testing.T) (*ConfigStage, *ConfigStage, *ConfigStage) {
	stage := &ConfigStage{
		t: t,
	}

	return stage, stage, stage
}

func (s *ConfigStage) And() *ConfigStage {
	return s
}

func (s *ConfigStage) AValidConfiguration() *ConfigStage {
	s.configRaw = []byte(`
pdAuthToken: abcdefghijklm
rotationStartHour: 08:00:00
currency: £
rotationPrices:
  - type: weekday
    price: 1
  - type: weekend
    price: 1
  - type: bankholiday
    price: 2
rotationUsers:
  - name: "User 1"
    holidaysCalendar: uk
    userId: ABCDEF1
  - name: "User 2"
    holidaysCalendar: uk
    userId: ABCDEF2
schedulesToIgnore:
  - SCHED_1
  - SCHED_2
  - SCHED_3
`)
	return s
}

func (s *ConfigStage) AMalformedConfiguration() *ConfigStage {
	s.configRaw = []byte(`
pdAuthToken: abcdefghijklm
	rotationStartHour: 08:00:00
  currency: £
		rotationPrices:
  	- type: weekday
    	price: 1
  - type: weekend
    price: 1
  	- type: bankholiday
    price: 2
	rotationUsers:
  - 	name: "User 1"
    holidaysCalendar: uk
    userId: ABCDEF1
  - name: "User 2"
    holidaysCalendar: uk
    userId: ABCDEF2
`)
	return s
}

func (s *ConfigStage) AValidConfigurationCorrectlyLoaded() *ConfigStage {
	s.AValidConfiguration().And().ItIsLoaded()
	assert.Nil(s.t, s.configError)
	return s
}

func (s *ConfigStage) ItIsLoaded() *ConfigStage {
	viper.SetConfigType("yaml")
	s.configError = viper.ReadConfig(bytes.NewBuffer(s.configRaw))
	if s.configError == nil {
		s.config = configuration.New()
		s.configUnmarshalError = viper.Unmarshal(s.config)
	}
	return s
}

func (s *ConfigStage) AnExistingPriceIsRequested() *ConfigStage {
	s.mapValue, s.mapError = s.config.FindPriceByDay("weekday")
	return s
}

func (s *ConfigStage) ANonExistingPriceIsRequested() *ConfigStage {
	s.mapValue, s.mapError = s.config.FindPriceByDay("wokday")
	return s
}

func (s *ConfigStage) AnExistingRotationInfoIsRequested() *ConfigStage {
	s.mapValue, s.mapError = s.config.FindRotationUserInfoByID("ABCDEF1")
	return s
}

func (s *ConfigStage) ANonExistingRotationInfoIsRequested() *ConfigStage {
	s.mapValue, s.mapError = s.config.FindRotationUserInfoByID("NONE")
	return s
}

func (s *ConfigStage) ValueIsFound() *ConfigStage {
	assert.Nil(s.t, s.mapError)
	assert.NotNil(s.t, s.mapValue)
	return s
}

func (s *ConfigStage) ValueIsNotFound() *ConfigStage {
	assert.NotNil(s.t, s.mapError)
	assert.Nil(s.t, s.mapValue)
	return s
}
func (s *ConfigStage) ConfigErrorIsCreated() *ConfigStage {
	assert.NotNil(s.t, s.configError)
	return s
}
