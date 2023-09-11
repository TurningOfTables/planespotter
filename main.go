package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"planespotter/helpers/config"
	"planespotter/helpers/formatters"
	"planespotter/helpers/save"
	"planespotter/helpers/types"
	"planespotter/helpers/ui"
	"time"

	"github.com/gen2brain/beeep"
	"golang.org/x/exp/slices"
)

const baseUrl = "opensky-network.org/api/states/all?"

var savePath = "save.json"

func main() {

	config, err := config.ParseConfig()
	if err != nil {
		log.Panic("Error reading config")
	}
	sa := calculateSearchArea(config.Position, config.SpotDistanceKm)
	err = save.CreateSaveIfNotExists(savePath)
	if err != nil {
		log.Panic("Error creating save file")
	}

	searchUrl, err := url.Parse(baseUrl)
	if err != nil {
		log.Panic(err)
	}

	queryParams := url.Values{
		"lamin": {sa.LaMin},
		"lomin": {sa.LoMin},
		"lamax": {sa.LaMax},
		"lomax": {sa.LoMax},
	}

	url := "https://" + config.ApiAuth.Username + ":" + config.ApiAuth.Password + "@" + searchUrl.String() + queryParams.Encode()

	go func() {
		for {
			log.Printf("Updating data...")
			planeInfos := updatePlanes(url)
			log.Printf("Received %v planes", len(planeInfos))
			notifyIfNew(planeInfos)
			log.Printf("Waiting %vs...", config.CheckFreqSeconds)
			time.Sleep(time.Duration(config.CheckFreqSeconds) * time.Second)

			log.Println("Finished waiting")
		}
	}()

	ui := ui.InitUi(savePath, config)
	ui.Run()

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

func updatePlanes(url string) []types.PlaneInfo {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		panic(errors.New("error getting response from server"))
	}

	remainingRequests := resp.Header.Get("X-Rate-Limit-Remaining")
	log.Printf("Remaining API requests today: %v\n", remainingRequests)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var r types.Result

	if err := json.Unmarshal(body, &r); err != nil {
		panic(err)
	}

	var planeInfos []types.PlaneInfo
	for _, plane := range r.PlaneResults {
		parsedResult := parseResult(plane)
		planeInfos = append(planeInfos, parsedResult)
	}

	return planeInfos
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
	saveData, err := save.GetSavedStats(savePath)
	if err != nil {
		log.Printf("Error getting saved stats: %v", err)
	}

	newPlanes := 0
	for _, p := range planeInfos {
		if slices.Contains(saveData.Callsigns, p.Icao24) {
			continue
		} else {
			newPlanes++
			save.SaveProgress(savePath, p)
			messageBody := fmt.Sprintf("%v \n ↑ %v → %v 🧭 %v \nTotal seen: %v", p.Callsign, p.Baro_Altitude, p.Velocity, p.True_Track, saveData.SeenCount+newPlanes)
			err = beeep.Notify("Plane Spotted!", messageBody, "assets/plane.png")
			if err != nil {
				log.Print(err)
			}

		}
	}
}
