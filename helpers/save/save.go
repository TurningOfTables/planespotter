package save

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"planespotter/helpers/types"

	"golang.org/x/exp/slices"
)

func SaveProgress(savePath string, p types.PlaneInfo) {
	err := CreateSaveIfNotExists(savePath)
	if err != nil {
		log.Printf("Error creating save: %v", err)
	}

	saveData, err := GetSave(savePath)
	if err != nil {
		log.Printf("Error getting current saved state")
	}

	if !slices.Contains(saveData.Callsigns, p.Callsign) {
		saveData.SeenCount += 1
		saveData.Callsigns = append(saveData.Callsigns, p.Callsign)
	}

	f, err := os.OpenFile(savePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error saving progress: %v", err)
	}
	defer f.Close()

	newSaveJson, err := json.MarshalIndent(saveData, "", "    ")
	if err != nil {
		log.Printf("Error saving progress: %v", err)
	}

	if err := os.Truncate(savePath, 0); err != nil {
		log.Printf("Error removing old save data")
	}

	err = os.WriteFile(savePath, newSaveJson, 0644)
	if err != nil {
		log.Printf("Error saving progress: %v", err)
	}
	f.Sync()
}

func SaveConfig(savePath string, saveConfig types.Config) {
	err := CreateSaveIfNotExists(savePath)
	if err != nil {
		log.Printf("Error creating save: %v", err)
	}

	currentSaveData, err := GetSave(savePath)
	if err != nil {
		log.Printf("Error getting current saved state")
	}

	var newSave = types.SaveData{
		Config:   saveConfig,
		Progress: currentSaveData.Progress,
	}

	newSaveJson, err := json.MarshalIndent(newSave, "", "    ")
	if err != nil {
		log.Printf("Error saving progress: %v", err)
	}

	if err := os.Truncate(savePath, 0); err != nil {
		log.Printf("Error removing old save data")
	}

	err = os.WriteFile(savePath, newSaveJson, 0644)
	if err != nil {
		log.Printf("Error saving progress: %v", err)
	}
}

func GetSave(savePath string) (types.SaveData, error) {
	var s types.SaveData
	saveFile, err := os.ReadFile(savePath)
	if err != nil {
		return s, err
	}

	json.Unmarshal(saveFile, &s)
	return s, nil
}

func CreateSaveIfNotExists(savePath string) error {
	_, err := os.Stat(savePath)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(savePath)
		if err != nil {
			return err
		}
		defer f.Close()

		var s types.SaveData

		// Set defaults
		s.Config.Position.Latitude = 40.730610
		s.Config.Position.Longitude = -73.935242
		s.Config.CheckFreqSeconds = 60
		s.Config.SpotDistanceKm = 20

		saveJson, err := json.MarshalIndent(s, "", "    ")
		if err != nil {
			return err
		}

		_, err = f.Write(saveJson)
		if err != nil {
			return err
		}
		f.Sync()
	}
	return nil
}
