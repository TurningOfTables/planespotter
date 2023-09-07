package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateSearchArea(t *testing.T) {

	invalidDataTests := []Position{
		{Latitude: 200, Longitude: 0},
		{Latitude: 0, Longitude: 200},
		{Latitude: 200, Longitude: 200},
	}

	for _, data := range invalidDataTests {
		if _, err := calculateSearchArea(data); err == nil {
			t.Errorf("calculateSearchArea(%+v) expected error, got nil", data)
		}

	}

	validDataTests := []Position{
		{Latitude: 40.730610, Longitude: -73.935242},
		{Latitude: 51.503399, Longitude: -0.119519},
		{Latitude: 22.302711, Longitude: 22.302711},
	}

	for _, data := range validDataTests {
		res, err := calculateSearchArea(data)
		if err != nil {
			t.Errorf("calculateSearchArea(%+v) expected nil, got error", data)
		}

		assert.IsType(t, SearchArea{}, res)
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
