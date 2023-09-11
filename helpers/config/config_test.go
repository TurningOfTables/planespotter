package config

import (
	"planespotter/helpers/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePosition(t *testing.T) {
	lat := 51.0
	long := 45.0

	res, err := parsePosition(lat, long)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, types.Position{Longitude: 45.000, Latitude: 51.000}, res)
}

func TestParsePositionErrors(t *testing.T) {

	tests := []struct {
		longitude     float64
		latitude      float64
		expectedError string
	}{
		{
			longitude:     999.000,
			latitude:      51.000,
			expectedError: "Longitude of 999 invalid. Must be between 90 and -90",
		},
		{
			longitude:     50.000,
			latitude:      222.000,
			expectedError: "Latitude of 222 invalid. Must be between 180 and -180",
		},
	}

	for _, test := range tests {

		_, err := parsePosition(test.latitude, test.longitude)
		if assert.Error(t, err) {
			assert.Equal(t, test.expectedError, err.Error())
		}
	}
}

func TestParseApiAuth(t *testing.T) {

	res, err := parseApiAuth("abc", "def")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, types.ApiAuth{Username: "abc", Password: "def"}, res)
}

func TestParseApiAuthErrors(t *testing.T) {

	tests := []struct {
		username      string
		password      string
		expectedError string
	}{
		{
			username:      "",
			password:      "def",
			expectedError: "error getting OPENSKY_USERNAME or OPENSKY_PASSWORD from env",
		},
		{
			username:      "abc",
			password:      "",
			expectedError: "error getting OPENSKY_USERNAME or OPENSKY_PASSWORD from env",
		},
	}

	for _, test := range tests {

		_, err := parseApiAuth(test.username, test.password)
		if assert.Error(t, err) {
			assert.Equal(t, test.expectedError, err.Error())
		}
	}
}
