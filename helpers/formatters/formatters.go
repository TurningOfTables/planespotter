// Package formatters provides functionality to format OpenSky API response fields
// into the correct output for planespotter notifications.

package formatters

import (
	"fmt"
	"math"
	"strings"
)

// FormatIcao24 takes an interface{} and returns the icao24 string if the underlying value is a string
// Otherwise it returns "N/A"
func FormatIcao24(icao24 interface{}) string {
	i, ok := icao24.(string)
	if !ok {
		return "N/A"
	}

	return i
}

// FormatCallsign takes an interface{} and returns the callsign string if the underlying value is a string
// Removes any extra spaces, as callsigns are padded to 8 chars by the OpenSky API
// Otherwise it returns "N/A"
func FormatCallsign(callsign interface{}) string {
	c, ok := callsign.(string)
	if !ok || c == "" {
		return "N/A"
	}

	c = strings.Replace(c, " ", "", -1)
	return c
}

// FormatBaroAltitude takes an interface{} and returns the baroaltitude string if the underlying value is a float64
// Converts from meters to feet, and appends ' ft' units
// Otherwise it returns "N/A"
func FormatBaroAltitude(baroAltitude interface{}) string {
	ba, ok := baroAltitude.(float64)
	if !ok {
		return "N/A"
	}

	ba *= 3.28084   // Convert meters to feet
	baFt := int(ba) // Convert to int
	return fmt.Sprintf("%v ft", baFt)
}

// FormatOnGround takes an interface{} and returns the onground string if the underlying value is a bool
// Otherwise it returns "N/A"
func FormatOnGround(onGround interface{}) string {
	og, ok := onGround.(bool)
	if !ok {
		return "N/A"
	}

	return fmt.Sprintf("%v", og)
}

// FormatVelocity takes an interface{} and returns the velocity string if the underlying value is a float64
// Converts from m/s to knots, and appends ' kts' units
// Otherwise returns "N/A"
func FormatVelocity(velocity interface{}) string {
	v, ok := velocity.(float64)
	if !ok {
		return "N/A"
	}

	v *= 1.94384 // convert m/s to knots
	vKt := int(v)
	return fmt.Sprintf("%v kts", vKt)
}

// FormatTrueTrack takes an interface{} and returns the trueTrack string if the underlying value is a float64
// Rounds to whole degrees
// Otherwise returns "N/A"
func FormatTrueTrack(trueTrack interface{}) string {
	tt, ok := trueTrack.(float64)
	if !ok {
		return "N/A"
	}

	ttRounded := int(tt)
	return fmt.Sprintf("%v°", ttRounded)
}

// KmToLatitude takes km int and converts it to decimal degrees of latitude
func KmToLatitude(km int) float64 {
	res := float64(km) / 111.1
	return res
}

// KmToLongitude takes km int and lat float64, and converts it to decimal degrees of longitude
func KmToLongitude(km int, lat float64) float64 {
	res := float64(km) / 110.320 * math.Cos(lat)
	return res
}
