package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"planespotter/helpers/formatters"
	"planespotter/helpers/save"
	"planespotter/helpers/types"
	"planespotter/helpers/ui"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/gen2brain/beeep"
	"golang.org/x/exp/slices"
)

const baseUrl = "opensky-network.org/api/states/all?"

const StartedText = "Spotting 🔭"
const StoppedText = "Stopped 🛑"
const ErrorText = "Error ⚠️"

var savePath = "save.json"
var pauseLoop = make(chan bool)
var started bool = false
var status = binding.NewString()

func main() {
	_, window := initApp()

	window.CenterOnScreen()
	window.ShowAndRun()
}

func initApp() (fyne.App, fyne.Window) {
	err := save.CreateSaveIfNotExists(savePath)
	if err != nil {
		log.Println("Error creating save file")
	}

	saveData, err := save.GetSave(savePath)
	if err != nil {
		log.Println("Error loading save file")
	}

	sa := calculateSearchArea(saveData.Config.Position, saveData.Config.SpotDistanceKm)
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

	app, window := initUi(url, savePath, saveData)
	return app, window
}

func initUi(url, savePath string, saveData types.SaveData) (fyne.App, fyne.Window) {
	icon, _ := fyne.LoadResourceFromPath("assets/plane.png")
	app := app.New()
	app.SetIcon(icon)
	window := ui.WindowSetup(app, icon)

	title := widget.NewLabelWithStyle("Configuration", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	settingsForm := ui.GenerateSettingsForm(savePath, saveData)

	startButton := widget.NewButton("Start", func() {
		if !started {
			log.Println("Spotting started")
			started = true
			status.Set(StartedText)
			go updateLoop(url, saveData)
		}
	})
	stopButton := widget.NewButton("Stop", func() {
		if started {
			log.Println("Spotting stopped")
			started = false
			status.Set(StoppedText)
			pauseLoop <- true
		}
	})

	if started {
		status.Set(StartedText)
	} else {
		status.Set(StoppedText)
	}
	statusLabel := widget.NewLabelWithData(status)

	window.SetContent(container.NewVBox(title, settingsForm, startButton, stopButton, statusLabel))
	return app, window
}

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

func calculateSearchArea(position types.Position, spotDistanceKm int) types.SearchArea {
	var sa types.SearchArea

	// Calculate offsets either side of position
	latSpotOffset := formatters.KmToLatitude(spotDistanceKm)
	longSpotOffset := formatters.KmToLongitude(spotDistanceKm, position.Latitude)

	// Finalise min and max long and lat by adding/subtracting offset
	sa.LaMax = fmt.Sprintf("%v", position.Latitude+latSpotOffset)
	sa.LaMin = fmt.Sprintf("%v", position.Latitude-latSpotOffset)
	sa.LoMax = fmt.Sprintf("%v", position.Longitude+longSpotOffset)
	sa.LoMin = fmt.Sprintf("%v", position.Longitude-longSpotOffset)
	return sa
}

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

func notifyIfNew(planeInfos []types.PlaneInfo) {
	saveData, err := save.GetSave(savePath)
	if err != nil {
		log.Printf("Error getting saved stats: %v", err)
	}

	newPlanes := 0
	for _, p := range planeInfos {
		if slices.Contains(saveData.Progress.Callsigns, p.Callsign) {
			continue
		} else {
			newPlanes++
			save.SaveProgress(savePath, p)
			messageBody := fmt.Sprintf("%v \n ↑ %v → %v 🧭 %v \nTotal seen: %v", p.Callsign, p.Baro_Altitude, p.Velocity, p.True_Track, saveData.SeenCount+newPlanes)
			err = beeep.Notify("Plane Spotted!", messageBody, "assets/plane.png")
			if err != nil {
				log.Printf("Error sending notification: %v", err)
			}

		}
	}
}
