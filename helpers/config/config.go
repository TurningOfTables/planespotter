package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"planespotter/helpers/types"
	"strconv"

	"github.com/joho/godotenv"
)

func ParseConfig() (types.Config, error) {
	var c types.Config
	err := godotenv.Load()
	if err != nil {
		return c, errors.New("error loading env")
	}

	lat, err := strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	if err != nil {
		return c, errors.New("error reading latitude from env")
	}

	long, err := strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)
	if err != nil {
		return c, errors.New("error reading longitude from env")
	}

	position, err := parsePosition(lat, long)
	if err != nil {
		return c, err
	}

	username := os.Getenv("OPENSKY_USERNAME")
	password := os.Getenv("OPENSKY_PASSWORD")
	apiAuth, err := parseApiAuth(username, password)
	if err != nil {
		return c, err
	}

	checkEverySeconds, err := strconv.Atoi(os.Getenv("CHECK_FREQ_SECONDS"))
	if err != nil {
		return c, errors.New("error reading check frequency from env")
	}

	spotDistanceKm, err := strconv.Atoi(os.Getenv("SPOT_DISTANCE_KM"))
	if err != nil {
		return c, errors.New("error reading spot distance from env")
	}

	c.Position = position
	c.ApiAuth = apiAuth
	c.CheckFreqSeconds = checkEverySeconds
	c.SpotDistanceKm = spotDistanceKm

	return c, nil
}

func SaveConfig(configPath string, c types.Config) error {

	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var newConfigOutput string

	newConfigOutput += "LATITUDE=" + fmt.Sprintf("%v", c.Position.Latitude) + "\n"
	newConfigOutput += "LONGITUDE=" + fmt.Sprintf("%v", c.Position.Longitude) + "\n"
	newConfigOutput += "OPENSKY_USERNAME=" + c.ApiAuth.Username + "\n"
	newConfigOutput += "OPENSKY_PASSWORD=" + c.ApiAuth.Password + "\n"
	newConfigOutput += "SPOT_DISTANCE_KM=" + fmt.Sprintf("%v", c.SpotDistanceKm) + "\n"
	newConfigOutput += "CHECK_FREQ_SECONDS=" + fmt.Sprintf("%v", c.CheckFreqSeconds) + "\n"

	if err := os.Truncate(configPath, 0); err != nil {
		log.Printf("Error removing old .env data")
	}

	_, err = f.WriteString(newConfigOutput)
	if err != nil {
		log.Fatal(err)
	}

	f.Sync()

	fmt.Println("Saved config")
	return nil
}

func parsePosition(lat, long float64) (types.Position, error) {
	var p types.Position

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

func parseApiAuth(username, password string) (types.ApiAuth, error) {
	var a types.ApiAuth

	if username == "" || password == "" {
		return a, errors.New("error getting OPENSKY_USERNAME or OPENSKY_PASSWORD from env")
	}

	a.Username = username
	a.Password = password
	return a, nil
}
