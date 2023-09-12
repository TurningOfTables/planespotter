package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatIcao24(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    "testicao",
			expected: "testicao",
		},
		{
			input:    555,
			expected: "N/A",
		},
		{
			input:    nil,
			expected: "N/A",
		},
	}

	for _, test := range tests {
		res := FormatIcao24(test.input)
		assert.Equal(t, test.expected, res)
	}
}

func TestFormatCallsign(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    "ABCD123",
			expected: "ABCD123",
		},
		{
			input:    "ABCD123 ",
			expected: "ABCD123",
		},
		{
			input:    nil,
			expected: "N/A",
		},
	}

	for _, test := range tests {
		res := FormatCallsign(test.input)
		assert.Equal(t, test.expected, res)
	}
}

func TestFormatBaroAltitude(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    1000.00,
			expected: "3280 ft",
		},
		{
			input:    "1000",
			expected: "N/A",
		},
		{
			input:    "",
			expected: "N/A",
		},
	}

	for _, test := range tests {
		res := FormatBaroAltitude(test.input)
		assert.Equal(t, test.expected, res)
	}
}

func TestFormatOnGround(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    true,
			expected: "true",
		},
		{
			input:    false,
			expected: "false",
		},
		{
			input:    "foo",
			expected: "N/A",
		},
	}

	for _, test := range tests {
		res := FormatOnGround(test.input)
		assert.Equal(t, test.expected, res)
	}
}

func TestFormatVelocity(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    100.00,
			expected: "194 kts",
		},
		{
			input:    "100",
			expected: "N/A",
		},
	}

	for _, test := range tests {
		res := FormatVelocity(test.input)
		assert.Equal(t, test.expected, res)
	}
}

func TestFormatTrueTrack(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    150.4,
			expected: "150°",
		},
		{
			input:    "fifty",
			expected: "N/A",
		},
	}

	for _, test := range tests {
		res := FormatTrueTrack(test.input)
		assert.Equal(t, test.expected, res)
	}
}

func TestKmToLatitude(t *testing.T) {
	tests := []struct {
		input    int
		expected float64
	}{
		{
			input:    50,
			expected: 0.45004500450045004,
		},
	}

	for _, test := range tests {
		res := KmToLatitude(test.input)
		assert.Equal(t, test.expected, res)
	}
}

func TestKmToLongitude(t *testing.T) {
	tests := []struct {
		inputKm  int
		inputLat float64
		expected float64
	}{
		{
			inputKm:  50,
			inputLat: 51.509865,
			expected: 0.14532643217451627,
		},
	}

	for _, test := range tests {
		res := KmToLongitude(test.inputKm, test.inputLat)
		assert.Equal(t, test.expected, res)
	}
}
