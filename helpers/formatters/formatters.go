package formatters

import (
	"fmt"
	"math"
	"strings"
)

func FormatIcao24(icao24 interface{}) string {
	i, ok := icao24.(string)
	if !ok {
		return "N/A"
	}

	return i
}

func FormatCallsign(callsign interface{}) string {
	c, ok := callsign.(string)
	if !ok || c == "" {
		return "N/A"
	}

	c = strings.Replace(c, " ", "", -1)
	return c
}

func FormatBaroAltitude(baroAltitude interface{}) string {
	ba, ok := baroAltitude.(float64)
	if !ok {
		return "N/A"
	}

	ba *= 3.28084   // Convert meters to feet
	baFt := int(ba) // Convert to int
	return fmt.Sprintf("%v ft", baFt)
}

func FormatOnGround(onGround interface{}) string {
	og, ok := onGround.(bool)
	if !ok {
		return "N/A"
	}

	return fmt.Sprintf("%v", og)
}

func FormatVelocity(velocity interface{}) string {
	v, ok := velocity.(float64)
	if !ok {
		return "N/A"
	}

	v *= 1.94384 // convert m/s to knots
	vKt := int(v)
	return fmt.Sprintf("%v kts", vKt)
}

func FormatTrueTrack(trueTrack interface{}) string {
	tt, ok := trueTrack.(float64)
	if !ok {
		return "N/A"
	}

	ttRounded := int(tt)
	return fmt.Sprintf("%v°", ttRounded)
}

func KmToLatitude(km int) float64 {
	res := float64(km) / 111.1
	return res
}

func KmToLongitude(km int, lat float64) float64 {
	res := float64(km) / 110.320 * math.Cos(lat)
	return res
}
