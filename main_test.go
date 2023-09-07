package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePosition(t *testing.T) {
	if err := os.Setenv("LONGITUDE", "45.000"); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("LATITUDE", "51.000"); err != nil {
		t.Error(err)
	}

	res, err := parsePosition()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, Position{Longitude: 45.000, Latitude: 51.000}, res)
}

func TestParsePositionErrors(t *testing.T) {
	if err := os.Setenv("LONGITUDE", "999.000"); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("LATITUDE", "51.000"); err != nil {
		t.Error(err)
	}

	_, err := parsePosition()
	if assert.Error(t, err) {
		assert.Equal(t, "Longitude of 999 invalid. Must be between 90 and -90", err.Error())
	}

	if err := os.Setenv("LONGITUDE", "50.000"); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("LATITUDE", "222.000"); err != nil {
		t.Error(err)
	}

	_, err = parsePosition()
	if assert.Error(t, err) {
		assert.Equal(t, "Latitude of 222 invalid. Must be between 180 and -180", err.Error())
	}
}

func TestParseApiAuth(t *testing.T) {
	if err := os.Setenv("OPENSKY_USERNAME", "abc"); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("OPENSKY_PASSWORD", "def"); err != nil {
		t.Error(err)
	}

	res, err := parseApiAuth()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, ApiAuth{username: "abc", password: "def"}, res)
}

func TestParseApiAuthErrors(t *testing.T) {
	if err := os.Setenv("OPENSKY_USERNAME", ""); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("OPENSKY_PASSWORD", "abc"); err != nil {
		t.Error(err)
	}

	_, err := parseApiAuth()
	if assert.Error(t, err) {
		assert.Equal(t, "Error getting OPENSKY_USERNAME or OPENSKY_PASSWORD from env", err.Error())
	}

	if err := os.Setenv("OPENSKY_USERNAME", "def"); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("OPENSKY_PASSWORD", ""); err != nil {
		t.Error(err)
	}

	_, err = parseApiAuth()
	if assert.Error(t, err) {
		assert.Equal(t, "Error getting OPENSKY_USERNAME or OPENSKY_PASSWORD from env", err.Error())
	}
}

func TestParseResult(t *testing.T) {
	// Good API response
	testApiResponse := []interface{}{
		"testicao", "testcallsign", "testorigincountry", 1234, 5678, 1111.2222, 3333.4444, 5555.66, true, 456.789, 123.456, 789.012, []int{1, 2, 3}, 987.654, "testsquawk", false, 1, 2,
	}
	expectedResult := PlaneInfo{
		Icao24:        "testicao",
		Callsign:      "testcallsign",
		Baro_Altitude: "18227 ft",
		On_Ground:     "true",
		Velocity:      "887 kts",
		True_Track:    "123°",
	}

	res := parseResult(testApiResponse)
	assert.IsType(t, PlaneInfo{}, res)
	assert.Equal(t, expectedResult, res, "parseResult(%+v) expected %+v, got %+v", testApiResponse, expectedResult, res)

	// Bad API response. Icao should be a string, but api returned an int!
	testBadApiResponse := []interface{}{
		1, "testcallsign", "testorigincountry", 1234, 5678, 1111.2222, 3333.4444, 5555.66, true, 456.789, 123.456, 789.012, []int{1, 2, 3}, 987.654, "testsquawk", false, 1, 2,
	}

	// N/A should be returned where a field isn't the expected type
	expectedBadResult := PlaneInfo{
		Icao24:        "N/A",
		Callsign:      "testcallsign",
		Baro_Altitude: "18227 ft",
		On_Ground:     "true",
		Velocity:      "887 kts",
		True_Track:    "123°",
	}

	res = parseResult(testBadApiResponse)
	assert.IsType(t, PlaneInfo{}, res)
	assert.Equal(t, expectedBadResult, res, "parseResult(%+v) expected %+v, got %+v", testApiResponse, expectedBadResult, res)

}

func TestCalculateSearchArea(t *testing.T) {
	p := Position{Longitude: 50.0, Latitude: 49.0}
	sa := SearchArea{laMax: "49.18001800180018", laMin: "48.81998199819982", loMax: "50.054494659852004", loMin: "49.945505340147996"}

	res := calculateSearchArea(p)

	assert.Less(t, res.laMin, res.laMax)
	assert.Less(t, res.loMin, res.loMax)
	assert.Equal(t, sa, res)
}

func TestUpdatePlanes(t *testing.T) {
	testApiPlane := []interface{}{
		"testicao", "testcallsign", "testorigincountry", 1234, 5678, 1111.2222, 3333.4444, 5555.66, true, 456.789, 123.456, 789.012, []int{1, 2, 3}, 987.654, "testsquawk", false, 1, 2,
	}

	var r = Result{PlaneResults: [][]interface{}{testApiPlane}}

	resBody, err := json.Marshal(r)
	if err != nil {
		t.Error("Error marshalling test data")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(resBody)
	}))
	defer server.Close()

	res := updatePlanes(server.URL)

	expectedResult := []PlaneInfo{
		{
			Icao24:        "testicao",
			Callsign:      "testcallsign",
			Baro_Altitude: "18227 ft",
			On_Ground:     "true",
			Velocity:      "887 kts",
			True_Track:    "123°"},
	}

	assert.Equal(t, expectedResult, res)
}

func TestUpdatePlanesErrors(t *testing.T) {
	failingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failingServer.Close()

	assert.PanicsWithError(t, `Error getting response from server`, func() { updatePlanes(failingServer.URL) })

	missingResBodyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	assert.PanicsWithError(t, `unexpected end of JSON input`, func() { updatePlanes(missingResBodyServer.URL) })

	badlyFormedResBodyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	}))

	assert.PanicsWithError(t, `unexpected end of JSON input`, func() { updatePlanes(badlyFormedResBodyServer.URL) })
}
