package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundCurrency(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{
			name:     "Round up from .666",
			input:    4.166666,
			expected: 4.17,
		},
		{
			name:     "Round down from .333",
			input:    3.333333,
			expected: 3.33,
		},
		{
			name:     "Already rounded value",
			input:    100.00,
			expected: 100.00,
		},
		{
			name:     "Round half to even (.165) - banker's rounding",
			input:    4.165,
			expected: 4.16,  // Go uses banker's rounding (round to even)
		},
		{
			name:     "Round half up (.5) alternate",
			input:    4.175,
			expected: 4.18,
		},
		{
			name:     "Small amount round up",
			input:    0.006,
			expected: 0.01,
		},
		{
			name:     "Small amount round down",
			input:    0.004,
			expected: 0.00,
		},
		{
			name:     "Large amount with decimals",
			input:    200.16666,
			expected: 200.17,
		},
		{
			name:     "Zero amount",
			input:    0.00,
			expected: 0.00,
		},
		{
			name:     "Negative amount (edge case)",
			input:    -4.166666,
			expected: -4.17,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := roundCurrency(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRoundCurrency_RealWorldScenarios(t *testing.T) {
	tests := []struct {
		name         string
		description  string
		hours        float32
		hourlyRate   float32
		expectedAmt  float32
	}{
		{
			name:         "30 minutes at £8.333/hr (24hr shift)",
			description:  "Half hour of a 24-hour £200 shift",
			hours:        0.5,
			hourlyRate:   8.333333,  // £200 ÷ 24 hours
			expectedAmt:  4.17,       // Should round to £4.17, not £4.166666
		},
		{
			name:         "1 hour at £8.333/hr (24hr shift)",
			description:  "One hour of a 24-hour £200 shift",
			hours:        1.0,
			hourlyRate:   8.333333,
			expectedAmt:  8.33,
		},
		{
			name:         "8 hours at £8.333/hr (24hr shift)",
			description:  "Eight hours of a 24-hour £200 shift",
			hours:        8.0,
			hourlyRate:   8.333333,
			expectedAmt:  66.67,
		},
		{
			name:         "15 hours at £6.666/hr (15hr shift)",
			description:  "Full 15-hour £100 shift",
			hours:        15.0,
			hourlyRate:   6.666666,  // £100 ÷ 15 hours
			expectedAmt:  100.00,     // Should equal exactly £100.00
		},
		{
			name:         "24 hours at £8.333/hr (24hr shift)",
			description:  "Full 24-hour £200 shift",
			hours:        24.0,
			hourlyRate:   8.333333,
			expectedAmt:  200.00,     // Should equal exactly £200.00
		},
		{
			name:         "30 minutes at £6.666/hr (15hr shift)",
			description:  "Half hour of a 15-hour £100 shift",
			hours:        0.5,
			hourlyRate:   6.666666,
			expectedAmt:  3.33,       // Should round to £3.33, not £3.333333
		},
		{
			name:         "Multiple 30-min intervals (5 hours)",
			description:  "10 intervals of 30-min at £8.333/hr",
			hours:        5.0,
			hourlyRate:   8.333333,
			expectedAmt:  41.67,
		},
		{
			name:         "Partial shift (7.5 hours)",
			description:  "15 intervals of 30-min at £8.333/hr",
			hours:        7.5,
			hourlyRate:   8.333333,
			expectedAmt:  62.50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate the amount and round it
			amount := tt.hours * tt.hourlyRate
			result := roundCurrency(amount)

			assert.Equal(t, tt.expectedAmt, result,
				"Expected %s: %.2f hours at £%.2f/hr = £%.2f (got £%.2f)",
				tt.description, tt.hours, tt.hourlyRate, tt.expectedAmt, result)
		})
	}
}

func TestRoundCurrency_FairnessCheck(t *testing.T) {
	// Verify that rounding doesn't systematically overpay or underpay
	// For a full 24-hour shift at £200, we should get close to £200

	hourlyRate := float32(200.0 / 24.0)  // £8.333333/hr

	// Simulate 48 half-hour intervals (24 hours)
	var total float32
	for i := 0; i < 48; i++ {
		total += 0.5 * hourlyRate
	}

	rounded := roundCurrency(total)

	// Should be exactly £200.00 (or within 1 penny due to float precision)
	assert.InDelta(t, 200.00, rounded, 0.01,
		"Full 24-hour shift should calculate to exactly £200.00")
}

func TestRoundCurrency_MultipleShifts(t *testing.T) {
	// Test scenario: User works multiple shifts across different days
	// This tests the summary calculation logic

	hourlyRate := float32(100.0 / 15.0)  // £6.666666/hr for 15-hour shift

	// Shift 1: 8 hours
	shift1 := roundCurrency(8.0 * hourlyRate)
	// Shift 2: 7.5 hours
	shift2 := roundCurrency(7.5 * hourlyRate)
	// Shift 3: 15 hours (full shift)
	shift3 := roundCurrency(15.0 * hourlyRate)

	// Total should be sum of rounded amounts
	total := roundCurrency(shift1 + shift2 + shift3)

	// Verify individual shifts are rounded
	assert.Equal(t, float32(53.33), shift1, "8 hours at £6.67/hr")
	assert.Equal(t, float32(50.00), shift2, "7.5 hours at £6.67/hr")
	assert.Equal(t, float32(100.00), shift3, "15 hours at £6.67/hr")

	// Total should be clean
	assert.Equal(t, float32(203.33), total, "Sum of multiple shifts")
}