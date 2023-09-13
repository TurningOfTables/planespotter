package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"planespotter/helpers/formatters"
	"planespotter/helpers/types"

	"golang.org/x/exp/slices"
)

func InitSaveData(savePath string) (string, types.SaveData) {
	saveData, err := GetSave(savePath)
	if err != nil {
		log.Println("Error loading save file")
	}

	sa := CalculateSearchArea(saveData.Config.Position, saveData.Config.SpotDistanceKm)
	searchUrl, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error parsing request URL")
	}

	queryParams := url.Values{
		"lamin": {sa.LaMin},
		"lomin": {sa.LoMin},
		"lamax": {sa.LaMax},
		"lomax": {sa.LoMax},
	}

	url := "https://" + saveData.ApiAuth.Username + ":" + saveData.ApiAuth.Password + "@" + searchUrl.String() + queryParams.Encode()

	return url, saveData
}

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

	SaveToFile(savePath, saveData)

}

func SaveConfig(savePath string, saveConfig types.Config) {
	err := CreateSaveIfNotExists(savePath)
	if err != nil {
		log.Printf("Error creating save: %v", err)
	}

	saveData, err := GetSave(savePath)
	if err != nil {
		log.Printf("Error getting current saved state")
	}

	var newSave = types.SaveData{
		Config:   saveConfig,
		Progress: saveData.Progress,
	}

	SaveToFile(savePath, newSave)
}

func SaveToFile(savePath string, saveData types.SaveData) {
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

func CalculateSearchArea(position types.Position, spotDistanceKm int) types.SearchArea {
	var sa types.SearchArea

	// Calculate offsets either side of position
	latSpotOffset := formatters.KmToLatitude(spotDistanceKm)
	longSpotOffset := formatters.KmToLongitude(spotDistanceKm, position.Latitude)

	// Finalise min and max long and lat by adding/subtracting offset
	sa.LaMax = fmt.Sprintf("%.4f", position.Latitude+latSpotOffset)
	sa.LaMin = fmt.Sprintf("%.4f", position.Latitude-latSpotOffset)
	sa.LoMax = fmt.Sprintf("%.4f", position.Longitude+longSpotOffset)
	sa.LoMin = fmt.Sprintf("%.4f", position.Longitude-longSpotOffset)
	return sa
}
