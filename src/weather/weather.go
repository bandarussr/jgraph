package weather

import (
	"encoding/json"
	"net/http"
	"net/url"
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

	// Build URL.
	u, err := url.Parse("https://nominatim.openstreetmap.org/reverse")
	if err != nil {
		panic("Error parsing URL: " + err.Error())
	}

	// Add parameters.
	q := u.Query()
	q.Set("lat", lat)
	q.Set("lon", long)
	q.Set("format", "json")
	q.Set("zoom", "10")
	u.RawQuery = q.Encode()

	// Make request.
	req, err := http.NewRequest("GET", u.String(), nil)
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

	// Read and parse.
	var nominatim nominatimResponse
	err = json.NewDecoder(res.Body).Decode(&nominatim)
	if err != nil {
		panic("Error parsing location data: " + err.Error())
	}

	return nominatim.Address.City
}
