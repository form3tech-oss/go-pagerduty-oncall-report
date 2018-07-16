package test

import (
	"testing"

	"github.com/rogersole/go-pagerduty-oncall-report/test/stages"
)

func TestFindExistingPriceByDay(t *testing.T) {
	given, when, then := stages.ConfigTest(t)

	given.
		A_valid_configuration_correctly_loaded()

	when.
		An_existing_price_is_requested()

	then.
		Value_is_found()
}

func TestFindNonExistingPriceByDay(t *testing.T) {
	given, when, then := stages.ConfigTest(t)

	given.
		A_valid_configuration_correctly_loaded()

	when.
		A_non_existing_price_is_requested()

	then.
		Value_is_not_found()
}

func TestFindExistingRotationUserInfoById(t *testing.T) {
	given, when, then := stages.ConfigTest(t)

	given.
		A_valid_configuration_correctly_loaded()

	when.
		An_existing_rotation_info_is_requested()

	then.
		Value_is_found()
}

func TestFindNonExistingRotationUserInfoById(t *testing.T) {
	given, when, then := stages.ConfigTest(t)

	given.
		A_valid_configuration_correctly_loaded()

	when.
		A_non_existing_rotation_info_is_requested()

	then.
		Value_is_not_found()
}

func TestConfigurationMalformed(t *testing.T) {
	given, when, then := stages.ConfigTest(t)

	given.
		A_malformed_configuration()

	when.
		It_is_loaded()

	then.
		Config_error_is_created()
}
