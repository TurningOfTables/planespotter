package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"planespotter/helpers"
	"strconv"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

const baseUrl = "opensky-network.org/api/states/all?"
const maxHistory = 5
const spotDistanceKm = 20
const checkEverySeconds = 60
const savePath = "save.json"

var seen []string

type SearchArea struct {
	laMin  string
	laMax  string
	loMin  string
	loMax  string
	height string
	width  string
	area   string
}

type ApiAuth struct {
	username string
	password string
}

type Position struct {
	Latitude  float64
	Longitude float64
}

type Result struct {
	PlaneResults [][]interface{} `json:"states"`
}

type PlaneInfo struct {
	Icao24        string
	Callsign      string
	Baro_Altitude string
	On_Ground     string
	Velocity      string
	True_Track    string
}

type SaveData struct {
	SeenCount int
	Callsigns []string
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading env vars.")
	}

	position, err := parsePosition()
	if err != nil {
		log.Panic(err)
	}

	sa := calculateSearchArea(position)

	auth, err := parseApiAuth()
	if err != nil {
		log.Panic(err)
	}

	searchUrl, err := url.Parse(baseUrl)
	if err != nil {
		log.Panic(err)
	}

	queryParams := url.Values{
		"lamin": {sa.laMin},
		"lomin": {sa.loMin},
		"lamax": {sa.laMax},
		"lomax": {sa.loMax},
	}

	url := "https://" + auth.username + ":" + auth.password + "@" + searchUrl.String() + queryParams.Encode()

	for {
		log.Printf("Updating data...")
		planeInfos := updatePlanes(url)
		notifyIfNew(planeInfos)
		log.Println("Waiting...")
		time.Sleep(checkEverySeconds * time.Second)
		log.Println("Finished waiting")
	}
}

func calculateSearchArea(position Position) SearchArea {
	var sa SearchArea

	// Calculate offsets either side of position
	latSpotOffset := helpers.KmToLatitude(spotDistanceKm)
	longSpotOffset := helpers.KmToLongitude(spotDistanceKm, position.Latitude)

	// Finalise min and max long and lat by adding/subtracting offset
	sa.laMax = fmt.Sprintf("%v", position.Latitude+latSpotOffset)
	sa.laMin = fmt.Sprintf("%v", position.Latitude-latSpotOffset)
	sa.loMax = fmt.Sprintf("%v", position.Longitude+longSpotOffset)
	sa.loMin = fmt.Sprintf("%v", position.Longitude-longSpotOffset)
	return sa
}

func parseResult(res []interface{}) PlaneInfo {
	var p PlaneInfo

	p.Icao24 = helpers.FormatIcao24(res[0])
	p.Callsign = helpers.FormatCallsign(res[1])
	p.Baro_Altitude = helpers.FormatBaroAltitude(res[7])
	p.On_Ground = helpers.FormatOnGround(res[8])
	p.Velocity = helpers.FormatVelocity(res[9])
	p.True_Track = helpers.FormatTrueTrack(res[10])

	return p
}

func parsePosition() (Position, error) {
	var p Position
	lat, err := strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	long, err := strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)
	if err != nil {
		return p, err
	}

	// Check long and lat are valid
	if lat > 180 || lat < -180 {
		error := fmt.Sprintf("Latitude of %v invalid. Must be between 180 and -180", lat)
		return p, errors.New(error)
	}
	if long > 90 || long < -90 {
		error := fmt.Sprintf("Longitude of %v invalid. Must be between 90 and -90", long)
		return p, errors.New(error)
	}

	p.Latitude = lat
	p.Longitude = long

	return p, nil
}

func parseApiAuth() (ApiAuth, error) {
	var a ApiAuth

	username := os.Getenv("OPENSKY_USERNAME")
	password := os.Getenv("OPENSKY_PASSWORD")

	if username == "" || password == "" {
		return a, errors.New("Error getting OPENSKY_USERNAME or OPENSKY_PASSWORD from env")
	}

	a.username = username
	a.password = password
	return a, nil
}

func updatePlanes(url string) []PlaneInfo {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		panic(errors.New("Error getting response from server"))
	}

	remainingRequests := resp.Header.Get("X-Rate-Limit-Remaining")
	log.Printf("Remaining API requests today: %v\n", remainingRequests)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var r Result

	if err := json.Unmarshal(body, &r); err != nil {
		panic(err)
	}

	var planeInfos []PlaneInfo
	for _, plane := range r.PlaneResults {
		parsedResult := parseResult(plane)
		planeInfos = append(planeInfos, parsedResult)
	}

	return planeInfos
}

func notifyIfNew(planeInfos []PlaneInfo) {
	for _, p := range planeInfos {
		if slices.Contains(seen, p.Icao24) {
			continue
		} else {
			seen = append(seen, p.Icao24)
			saveProgress(p)
			stats, err := getSavedStats()
			if err != nil {
				log.Print("Error getting saved stats")
			}
			messageBody := fmt.Sprintf("%v \n ↑ %v → %v 🧭 %v \nTotal seen: %v", p.Callsign, p.Baro_Altitude, p.Velocity, p.True_Track, stats.SeenCount)
			err = beeep.Notify("Plane Spotted!", messageBody, "assets/plane.png")
			if err != nil {
				log.Print(err)
			}

		}
	}

	// trim seen history so it doesn't get too big
	if len(seen) > maxHistory {
		numOverMax := len(seen) - maxHistory
		seen = slices.Delete(seen, 0, numOverMax)
	}
}

func saveProgress(p PlaneInfo) {
	err := createSaveIfNotExists()
	if err != nil {
		log.Printf("Error creating save: %v", err)
	}

	saveFile, err := os.ReadFile(savePath)
	if err != nil {
		log.Printf("Error reading save: %v", err)
	}
	var s SaveData

	json.Unmarshal(saveFile, &s)

	if !slices.Contains(s.Callsigns, p.Callsign) {
		s.SeenCount += 1
		s.Callsigns = append(s.Callsigns, p.Callsign)
	}

	f, err := os.OpenFile(savePath, os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error saving progress: %v", err)
	}
	defer f.Close()

	newSaveJson, err := json.Marshal(s)
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

func getSavedStats() (SaveData, error) {
	var s SaveData
	saveFile, err := os.ReadFile(savePath)
	if err != nil {
		return s, err
	}

	json.Unmarshal(saveFile, &s)
	return s, nil
}

func createSaveIfNotExists() error {
	_, err := os.Stat(savePath)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(savePath)
		if err != nil {
			return err
		}
		defer f.Close()

		var s = SaveData{SeenCount: 0, Callsigns: []string{}}
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
