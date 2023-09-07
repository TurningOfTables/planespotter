package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testSavePath = "test_save.json"

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

	tests := []struct {
		longitude     string
		latitude      string
		expectedError string
	}{
		{
			longitude:     "999.000",
			latitude:      "51.000",
			expectedError: "Longitude of 999 invalid. Must be between 90 and -90",
		},
		{
			longitude:     "50.000",
			latitude:      "222.000",
			expectedError: "Latitude of 222 invalid. Must be between 180 and -180",
		},
		{
			longitude:     "foo",
			latitude:      "51.000",
			expectedError: `strconv.ParseFloat: parsing "foo": invalid syntax`,
		},
		{
			longitude:     "48.000",
			latitude:      "bar",
			expectedError: `strconv.ParseFloat: parsing "bar": invalid syntax`,
		},
	}

	for _, test := range tests {
		os.Setenv("LONGITUDE", test.longitude)
		os.Setenv("LATITUDE", test.latitude)

		_, err := parsePosition()
		if assert.Error(t, err) {
			assert.Equal(t, test.expectedError, err.Error())
		}
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
		os.Setenv("OPENSKY_USERNAME", test.username)
		os.Setenv("OPENSKY_PASSWORD", test.password)

		_, err := parseApiAuth()
		if assert.Error(t, err) {
			assert.Equal(t, test.expectedError, err.Error())
		}
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

	assert.PanicsWithError(t, `error getting response from server`, func() { updatePlanes(failingServer.URL) })

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

func TestCreateSaveIfNotExists(t *testing.T) {
	_, err := os.Stat(testSavePath)
	if !errors.Is(err, os.ErrNotExist) {
		os.Remove(testSavePath)
	}

	err = createSaveIfNotExists(testSavePath)
	if err != nil {
		t.Error(err)
	}

	f, err := os.ReadFile(testSavePath)
	if err != nil {
		t.Error(err)
	}

	assert.FileExists(t, testSavePath)
	var s SaveData

	json.Unmarshal(f, &s)
	assert.Equal(t, 0, s.SeenCount)
	assert.IsType(t, []string{}, s.Callsigns)

	err = os.Remove(testSavePath)
	if err != nil {
		t.Error(err)
	}

}

func TestSaveProgress(t *testing.T) {
	err := createSaveIfNotExists(testSavePath)
	if err != nil {
		t.Error(err)
	}

	p := PlaneInfo{
		Icao24:        "testicao",
		Callsign:      "testcallsign",
		Baro_Altitude: "18227 ft",
		On_Ground:     "true",
		Velocity:      "887 kts",
		True_Track:    "123°",
	}

	saveProgress(testSavePath, p)

	var s SaveData
	saveFile, err := os.ReadFile(testSavePath)
	if err != nil {
		t.Error(err)
	}

	json.Unmarshal(saveFile, &s)

	expectedSave := SaveData{SeenCount: 1, Callsigns: []string{"testcallsign"}}
	assert.Equal(t, expectedSave, s)

	err = os.Remove(testSavePath)
	if err != nil {
		t.Error(err)
	}
}

func TestGetSavedStats(t *testing.T) {
	err := createSaveIfNotExists(testSavePath)
	if err != nil {
		t.Error(err)
	}

	p := PlaneInfo{
		Icao24:        "testicao",
		Callsign:      "testcallsign",
		Baro_Altitude: "18227 ft",
		On_Ground:     "true",
		Velocity:      "887 kts",
		True_Track:    "123°",
	}

	saveProgress(testSavePath, p)

	s, err := getSavedStats(testSavePath)
	if err != nil {
		t.Error(err)
	}

	expectedSave := SaveData{SeenCount: 1, Callsigns: []string{"testcallsign"}}
	assert.Equal(t, expectedSave, s)

	err = os.Remove(testSavePath)
	if err != nil {
		t.Error(err)
	}
}
