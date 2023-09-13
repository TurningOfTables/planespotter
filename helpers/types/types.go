// Package types defines various structs to be shared across the app

package types

type SearchArea struct {
	LaMin  string
	LaMax  string
	LoMin  string
	LoMax  string
	Height string
	Width  string
	Area   string
}

type ApiAuth struct {
	Username string
	Password string
}

type Position struct {
	Latitude  float64
	Longitude float64
}

type Result struct {
	PlaneResults [][]interface{} `json:"states"`
}

type Config struct {
	Position         Position
	ApiAuth          ApiAuth
	SpotDistanceKm   int
	CheckFreqSeconds int
}

type Progress struct {
	SeenCount int
	Callsigns []string
}

type SaveData struct {
	Config
	Progress
}

type PlaneInfo struct {
	Icao24        string
	Callsign      string
	Baro_Altitude string
	On_Ground     string
	Velocity      string
	True_Track    string
}
