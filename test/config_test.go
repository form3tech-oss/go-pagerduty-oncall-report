package test

import (
	"testing"

	"github.com/rogersole/go-pagerduty-oncall-report/test/stages"
)

func TestFindExistingPriceByDay(t *testing.T) {
	given, when, then := stages.ConfigTest(t)

	given.
		AValidConfigurationCorrectlyLoaded()

	when.
		AnExistingPriceIsRequested()

	then.
		ValueIsFound()
}

func TestFindNonExistingPriceByDay(t *testing.T) {
	given, when, then := stages.ConfigTest(t)

	given.
		AValidConfigurationCorrectlyLoaded()

	when.
		ANonExistingPriceIsRequested()

	then.
		ValueIsNotFound()
}

func TestFindExistingRotationUserInfoById(t *testing.T) {
	given, when, then := stages.ConfigTest(t)

	given.
		AValidConfigurationCorrectlyLoaded()

	when.
		AnExistingRotationInfoIsRequested()

	then.
		ValueIsFound()
}

func TestFindNonExistingRotationUserInfoById(t *testing.T) {
	given, when, then := stages.ConfigTest(t)

	given.
		AValidConfigurationCorrectlyLoaded()

	when.
		ANonExistingRotationInfoIsRequested()

	then.
		ValueIsNotFound()
}

func TestConfigurationMalformed(t *testing.T) {
	given, when, then := stages.ConfigTest(t)

	given.
		AMalformedConfiguration()

	when.
		ItIsLoaded()

	then.
		ConfigErrorIsCreated()
}
