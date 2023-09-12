package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"planespotter/helpers/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseResult(t *testing.T) {
	// Good API response
	testApiResponse := []interface{}{
		"testicao", "testcallsign", "testorigincountry", 1234, 5678, 1111.2222, 3333.4444, 5555.66, true, 456.789, 123.456, 789.012, []int{1, 2, 3}, 987.654, "testsquawk", false, 1, 2,
	}
	expectedResult := types.PlaneInfo{
		Icao24:        "testicao",
		Callsign:      "testcallsign",
		Baro_Altitude: "18227 ft",
		On_Ground:     "true",
		Velocity:      "887 kts",
		True_Track:    "123°",
	}

	res := parseResult(testApiResponse)
	assert.IsType(t, types.PlaneInfo{}, res)
	assert.Equal(t, expectedResult, res, "parseResult(%+v) expected %+v, got %+v", testApiResponse, expectedResult, res)

	// Bad API response. Icao should be a string, but api returned an int!
	testBadApiResponse := []interface{}{
		1, "testcallsign", "testorigincountry", 1234, 5678, 1111.2222, 3333.4444, 5555.66, true, 456.789, 123.456, 789.012, []int{1, 2, 3}, 987.654, "testsquawk", false, 1, 2,
	}

	// N/A should be returned where a field isn't the expected type
	expectedBadResult := types.PlaneInfo{
		Icao24:        "N/A",
		Callsign:      "testcallsign",
		Baro_Altitude: "18227 ft",
		On_Ground:     "true",
		Velocity:      "887 kts",
		True_Track:    "123°",
	}

	res = parseResult(testBadApiResponse)
	assert.IsType(t, types.PlaneInfo{}, res)
	assert.Equal(t, expectedBadResult, res, "parseResult(%+v) expected %+v, got %+v", testApiResponse, expectedBadResult, res)

}

func TestUpdatePlanes(t *testing.T) {
	testApiPlane := []interface{}{
		"testicao", "testcallsign", "testorigincountry", 1234, 5678, 1111.2222, 3333.4444, 5555.66, true, 456.789, 123.456, 789.012, []int{1, 2, 3}, 987.654, "testsquawk", false, 1, 2,
	}

	var r = types.Result{PlaneResults: [][]interface{}{testApiPlane}}

	resBody, err := json.Marshal(r)
	if err != nil {
		t.Error("Error marshalling test data")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(resBody)
	}))
	defer server.Close()

	res, err := updatePlanes(server.URL)
	if err != nil {
		t.Error(err)
	}

	expectedResult := []types.PlaneInfo{
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

	_, err := updatePlanes(failingServer.URL)
	assert.Error(t, err)

	missingResBodyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	_, err = updatePlanes(missingResBodyServer.URL)
	assert.Error(t, err)

	badlyFormedResBodyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	}))

	_, err = updatePlanes(badlyFormedResBodyServer.URL)
	assert.Error(t, err)
}
