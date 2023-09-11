package save

import (
	"encoding/json"
	"errors"
	"os"
	"planespotter/helpers/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testSavePath = "test_save.json"

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

	expectedSave := types.SaveData{SeenCount: 1, Callsigns: []string{"testcallsign"}}
	assert.Equal(t, expectedSave, s)

	err = os.Remove(testSavePath)
	if err != nil {
		t.Error(err)
	}
}

func TestGetSavedStats(t *testing.T) {
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

	s, err := GetSavedStats(testSavePath)
	if err != nil {
		t.Error(err)
	}

	expectedSave := types.SaveData{SeenCount: 1, Callsigns: []string{"testcallsign"}}
	assert.Equal(t, expectedSave, s)

	err = os.Remove(testSavePath)
	if err != nil {
		t.Error(err)
	}
}
