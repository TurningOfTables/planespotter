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

type SaveData struct {
	SeenCount int
	Callsigns []string
}

type PlaneInfo struct {
	Icao24        string
	Callsign      string
	Baro_Altitude string
	On_Ground     string
	Velocity      string
	True_Track    string
}

type Config struct {
	Position         Position
	ApiAuth          ApiAuth
	SpotDistanceKm   int
	CheckFreqSeconds int
}
