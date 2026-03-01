package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
)

type openMeteoResponse struct {
	Daily struct {
		Time                        []int     `json:"time"`
		WeatherCode                 []int     `json:"weather_code"`
		Temperature2mMin            []float32 `json:"temperature_2m_min"`
		Temperature2mMax            []float32 `json:"temperature_2m_max"`
		PrecipitationProbabilityMax []int     `json:"precipitation_probability_max"`
		WindSpeed10mMax             []float32 `json:"wind_speed_10m_max"`
		WindDirection10mDominant    []int     `json:"wind_direction_10m_dominant"`
	} `json:"daily"`
}

func getWeatherFromOpenMeteo(lat, long string) *openMeteoResponse {
	// Limit API requests during development.
	if dev {
		data, err := os.ReadFile(".devdata.json")
		if err == nil {
			var weather openMeteoResponse
			err = json.Unmarshal(data, &weather)
			if err == nil {
				return &weather
			}
		}
	}

	// Build URL.
	u, err := url.Parse("https://api.open-meteo.com/v1/forecast")
	if err != nil {
		panic("Error parsing URL: " + err.Error())
	}

	// Add parameters.
	q := u.Query()
	q.Set("latitude", lat)
	q.Set("longitude", long)
	q.Set("timezone", "auto")
	q.Set("wind_speed_unit", "mph")
	q.Set("temperature_unit", "fahrenheit")
	q.Set("precipitation_unit", "inch")
	q.Set("timeformat", "unixtime")

	daily := []string{
		"weather_code",
		"temperature_2m_min",
		"temperature_2m_max",
		"precipitation_probability_max",
		"wind_speed_10m_max",
		"wind_direction_10m_dominant",
	}
	q.Set("daily", strings.Join(daily, ","))
	u.RawQuery = q.Encode()

	// Make request.
	res, err := http.Get(u.String())
	if err != nil {
		panic("Error getting weather data.")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic("Error getting weather data: " + res.Status)
	}

	// Read and parse.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic("Error reading weather data.")
	}

	var weather openMeteoResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic("Error pasing weather data: " + err.Error())
	}

	// Save data for development.
	if dev {
		err = os.WriteFile(".devdata.json", body, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing weather data to file: %s\n", err.Error())
		}
	}

	return &weather
}

func convertIsoTimesToStrings(times []int) []string {
	today := time.Now()
	timeStr := make([]string, len(times))
	
	for i, k := range times {
		t := time.Unix(int64(k), 0)
		
		if t.Day() == today.Day() && t.Month() == today.Month() && t.Year() == today.Year() {
			timeStr[i] = "Today"
			continue
		}
		
		timeStr[i] = t.Weekday().String()[:3]
	}

	return timeStr
}

func (o *openMeteoResponse) toWeather() *Weather {
	w := Weather{
		Keys:            convertIsoTimesToStrings(o.Daily.Time),
		MinTemperature:  slices.Min(o.Daily.Temperature2mMin),
		MaxTemperature:  slices.Max(o.Daily.Temperature2mMax),
		TemperatureHigh: make(map[string]float32),
		TemperatureLow:  make(map[string]float32),
		Condition:       make(map[string]Condition),
	}

	for i, key := range w.Keys {
		w.TemperatureHigh[key] = o.Daily.Temperature2mMax[i]
		w.TemperatureLow[key] = o.Daily.Temperature2mMin[i]
		w.Condition[key] = Condition(o.Daily.WeatherCode[i])
	}

	return &w
}
