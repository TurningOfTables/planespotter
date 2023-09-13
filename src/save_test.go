package main

import (
	"encoding/json"
	"errors"
	"os"
	"planespotter/helpers/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSaveIfNotExists(t *testing.T) {
	_, err := os.Stat(testSavePath)
	if !errors.Is(err, os.ErrNotExist) {
		os.Remove(testSavePath)
	}

	err = CreateSaveIfNotExists(testSavePath)
	if err != nil {
		t.Error(err)
	}

	f, err := os.ReadFile(testSavePath)
	if err != nil {
		t.Error(err)
	}

	assert.FileExists(t, testSavePath)
	var s types.SaveData

	json.Unmarshal(f, &s)
	assert.Equal(t, 0, s.SeenCount)
	assert.IsType(t, []string{}, s.Callsigns)

	err = os.Remove(testSavePath)
	if err != nil {
		t.Error(err)
	}

}

func TestSaveProgress(t *testing.T) {
	err := CreateSaveIfNotExists(testSavePath)
	if err != nil {
		t.Error(err)
	}

	p := types.PlaneInfo{
		Icao24:        "testicao",
		Callsign:      "testcallsign",
		Baro_Altitude: "18227 ft",
		On_Ground:     "true",
		Velocity:      "887 kts",
		True_Track:    "123°",
	}

	SaveProgress(testSavePath, p)

	var s types.SaveData
	saveFile, err := os.ReadFile(testSavePath)
	if err != nil {
		t.Error(err)
	}

	json.Unmarshal(saveFile, &s)

	expectedSave := types.SaveData{Config: types.Config{Position: types.Position{Latitude: 40.73061, Longitude: -73.935242}, ApiAuth: types.ApiAuth{Username: "", Password: ""}, SpotDistanceKm: 20, CheckFreqSeconds: 60}, Progress: types.Progress{SeenCount: 1, Callsigns: []string{"testcallsign"}}}
	assert.Equal(t, expectedSave, s)

	err = os.Remove(testSavePath)
	if err != nil {
		t.Error(err)
	}
}

func TestSaveConfig(t *testing.T) {
	err := CreateSaveIfNotExists(testSavePath)
	if err != nil {
		t.Error(err)
	}

	c := types.Config{
		Position:         types.Position{Latitude: 48.00, Longitude: 2},
		ApiAuth:          types.ApiAuth{Username: "blah", Password: "fluff"},
		CheckFreqSeconds: 9,
		SpotDistanceKm:   8,
	}

	SaveConfig(testSavePath, c)

	save, err := GetSave(testSavePath)
	if err != nil {
		t.Error(err)
	}

	expectedSave := types.SaveData(types.SaveData{Config: types.Config{Position: types.Position{Latitude: 48, Longitude: 2}, ApiAuth: types.ApiAuth{Username: "blah", Password: "fluff"}, SpotDistanceKm: 8, CheckFreqSeconds: 9}, Progress: types.Progress{SeenCount: 0, Callsigns: []string(nil)}})

	assert.Equal(t, expectedSave, save)

	err = os.Remove(testSavePath)
	if err != nil {
		t.Error(err)
	}
}

func TestGetSave(t *testing.T) {
	err := CreateSaveIfNotExists(testSavePath)
	if err != nil {
		t.Error(err)
	}

	p := types.PlaneInfo{
		Icao24:        "testicao",
		Callsign:      "testcallsign",
		Baro_Altitude: "18227 ft",
		On_Ground:     "true",
		Velocity:      "887 kts",
		True_Track:    "123°",
	}

	SaveProgress(testSavePath, p)

	s, err := GetSave(testSavePath)
	if err != nil {
		t.Error(err)
	}

	expectedSave := types.SaveData{Config: types.Config{Position: types.Position{Latitude: 40.73061, Longitude: -73.935242}, ApiAuth: types.ApiAuth{Username: "", Password: ""}, SpotDistanceKm: 20, CheckFreqSeconds: 60}, Progress: types.Progress{SeenCount: 1, Callsigns: []string{"testcallsign"}}}
	assert.Equal(t, expectedSave, s)

	err = os.Remove(testSavePath)
	if err != nil {
		t.Error(err)
	}
}

func TestCalculateSearchArea(t *testing.T) {
	p := types.Position{Longitude: 50.0, Latitude: 49.0}
	sa := types.SearchArea{LaMin: "48.9100", LaMax: "49.0900", LoMin: "49.9728", LoMax: "50.0272", Height: "", Width: "", Area: ""}

	res := CalculateSearchArea(p, 10)

	assert.Less(t, res.LaMin, res.LaMax)
	assert.Less(t, res.LoMin, res.LoMax)
	assert.Equal(t, sa, res)
}
