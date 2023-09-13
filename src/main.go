package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"planespotter/helpers/formatters"
	"planespotter/helpers/types"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/gen2brain/beeep"
	"golang.org/x/exp/slices"
)

const baseUrl = "opensky-network.org/api/states/all?"

const StartedText = "Spotting 🔭"
const StoppedText = "Stopped 🛑"
const ErrorText = "Error ⚠️"

var savePath = "save.json"
var testSavePath = "test_save.json"
var pauseLoop = make(chan bool)
var started bool = false
var status = binding.NewString()

// main creates a new save if required
// Loads save data
// Starts the UI
func main() {
	err := CreateSaveIfNotExists(savePath)
	if err != nil {
		log.Println("Error creating save file")
	}

	url, saveData := InitSaveData(savePath)
	_, window := InitUi(url, savePath, saveData)

	window.CenterOnScreen()
	window.ShowAndRun()
}

// startUpdateLoop takes the API url and saveData, and begins the updateLoop to check for new planes and notify the user
func startUpdateLoop(url string, saveData types.SaveData) {
	log.Println("Spotting started")
	started = true
	status.Set(StartedText)
	go updateLoop(url, saveData)
}

// stopUpdateLoop stops checking for new planes and notifying the user
func stopUpdateLoop() {
	log.Println("Spotting stopped")
	started = false
	status.Set(StoppedText)
	pauseLoop <- true
}

// updateLoop takes the API url and saveData, and triggers a check for new planes
// It sends the results to notifyIfNew to send notifications
// It checks based on the check frequency specified in the config
// It stops when it receives on the pauseLoop channel
func updateLoop(url string, saveData types.SaveData) {
	timeSinceCheck := saveData.CheckFreqSeconds
	for range time.Tick(time.Second) {
		timeSinceCheck += 1
		select {
		case <-pauseLoop:
			return
		default:
			if timeSinceCheck >= saveData.CheckFreqSeconds {
				planeInfos, err := updatePlanes(url)
				if err != nil {
					log.Printf("Error updating planes: %v", err)
					status.Set(ErrorText + " - " + err.Error())
					started = false
					pauseLoop <- true
				}
				log.Printf("Received %v planes", len(planeInfos))
				notifyIfNew(planeInfos)
				timeSinceCheck = 0
			}
		}
	}
}

// updatePlanes takes the API url and returns a slice of PlaneInfo
// It also logs out the remaining API requests to help monitor API usage
// Returns an error if the request to the API or parsing the response fails
func updatePlanes(url string) ([]types.PlaneInfo, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		errorString := fmt.Sprintf("error getting response from server: status %v | error %v", resp.StatusCode, err)
		return []types.PlaneInfo{}, errors.New(errorString)
	}

	remainingRequests := resp.Header.Get("X-Rate-Limit-Remaining")
	log.Printf("Remaining API requests today: %v\n", remainingRequests)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorString := fmt.Sprintf("error reading response from server: status %v | error %v", resp.StatusCode, err)
		return []types.PlaneInfo{}, errors.New(errorString)
	}
	defer resp.Body.Close()

	var r types.Result

	if err := json.Unmarshal(body, &r); err != nil {
		errorString := fmt.Sprintf("error parsing response from server: status %v | error %v", resp.StatusCode, err)
		return []types.PlaneInfo{}, errors.New(errorString)
	}

	var planeInfos []types.PlaneInfo
	for _, plane := range r.PlaneResults {
		parsedResult := parseResult(plane)
		planeInfos = append(planeInfos, parsedResult)
	}

	return planeInfos, nil
}

// parseResult takes an individual plane response from the API and assigns it to a PlaneInfo object's fields, then returns the PlaneInfo
// Based on the position in the response as the API doesn't return key/value, just an array
func parseResult(res []interface{}) types.PlaneInfo {
	var p types.PlaneInfo

	p.Icao24 = formatters.FormatIcao24(res[0])
	p.Callsign = formatters.FormatCallsign(res[1])
	p.Baro_Altitude = formatters.FormatBaroAltitude(res[7])
	p.On_Ground = formatters.FormatOnGround(res[8])
	p.Velocity = formatters.FormatVelocity(res[9])
	p.True_Track = formatters.FormatTrueTrack(res[10])

	return p
}

// notifyIfNew takes the slice of PlaneInfo
// For each one, if it has been seen before then it continues without doing anything
// It it has not been seen before it sends a notification to the user, and updates the save progress
func notifyIfNew(planeInfos []types.PlaneInfo) {
	saveData, err := GetSave(savePath)
	if err != nil {
		log.Printf("Error getting saved stats: %v", err)
	}

	newPlanes := 0
	for _, p := range planeInfos {
		if slices.Contains(saveData.Progress.Callsigns, p.Callsign) {
			continue
		} else {
			newPlanes++
			SaveProgress(savePath, p)
			messageBody := fmt.Sprintf("%v \n ↑ %v → %v 🧭 %v \nTotal seen: %v", p.Callsign, p.Baro_Altitude, p.Velocity, p.True_Track, saveData.SeenCount+newPlanes)
			err = beeep.Notify("Plane Spotted!", messageBody, "assets/plane.png")
			if err != nil {
				log.Printf("Error sending notification: %v", err)
			}

		}
	}
}
