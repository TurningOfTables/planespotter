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

	saveData, err := GetSavedStats(savePath)
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

	newSaveJson, err := json.Marshal(saveData)
	if err != nil {
		log.Printf("Error saving progress: %v", err)
	}

	if err := os.Truncate(savePath, 0); err != nil {
		log.Printf("Error removing old save data")
	}

	_, err = f.Write(newSaveJson)
	if err != nil {
		log.Printf("Error saving progress: %v", err)
	}
	f.Sync()
}

func GetSavedStats(savePath string) (types.SaveData, error) {
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

		var s = types.SaveData{SeenCount: 0, Callsigns: []string{}}
		saveJson, err := json.Marshal(s)
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
