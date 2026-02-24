package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const dev = true

type Weather struct {
	Location        string
	Keys            []string
	TemperatureHigh map[string]float32
	TemperatureLow  map[string]float32
	Condition       map[string]Condition
}

type Condition int

const (
	ConditionSun Condition = iota
	ConditionCloud
	ConditionRain
	ConditionSnow
	ConditionThunderstorm
)

type Forecast int

const (
	ForecastDaily Forecast = iota
	ForecastHourly
)

func New(lat, long string, f Forecast) *Weather {
	w := getWeatherFromOpenMeteo(lat, long).toWeather()
	w.Location = getLocationName(lat, long)
	return w
}

func getLocationName(lat, long string) string {
	type nominatimResponse struct {
		Address struct {
			City string `json:"city"`
		} `json:"address"`
	}

	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s&format=json&zoom=10", lat, long)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Error creating request: " + err.Error())
	}

	// User-Agent header (required by Nominatim)
	req.Header.Set("User-Agent", "JgraphLab/1.0 (github.com/bandarussr/jgraph)")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic("Error getting location data.")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic("Error getting location data: " + res.Status)
	}

	var nominatim nominatimResponse
	err = json.NewDecoder(res.Body).Decode(&nominatim)
	if err != nil {
		panic("Error parsing location data: " + err.Error())
	}

	return nominatim.Address.City
}
