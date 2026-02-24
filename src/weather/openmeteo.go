package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type openMeteoResponse struct {
	// Hourly struct {
	// 	Time                     []string  `json:"time"`
	// 	Temperature2m            []float32 `json:"temperature_2m"`
	// 	RelativeHumidity2m       []int     `json:"relative_humidity_2m"`
	// 	PrecipitationProbability []int     `json:"precipitation_probability"`
	// 	WeatherCode              []int     `json:"weather_code"`
	// 	WindSpeed10m             []float32 `json:"wind_speed_10m"`
	// 	WindDirection10m         []int     `json:"wind_direction_10m"`
	// } `json:"hourly"`
	Daily struct {
		Time                        []string  `json:"time"`
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
	
	daily := []string{
		"weather_code",
		"temperature_2m_min",
		"temperature_2m_max",
		"precipitation_probability_max",
		"wind_speed_10m_max",
		"wind_direction_10m_dominant",
	}
	q.Set("daily", strings.Join(daily, ","))
	
	// hourly := []string{
	// 	"temperature_2m",
	// 	"relative_humidity_2m",
	// 	"precipitation_probability",
	// 	"weather_code",
	// 	"wind_speed_10m",
	// 	"wind_direction_10m",
	// }
	// q.Set("hourly", strings.Join(hourly, ","))
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

func convertIsoTimesToStrings(times []string) []string {
	panic("Not implemented")
}

func (o *openMeteoResponse) toWeather() *Weather {
	w := Weather{
		Keys:            o.Daily.Time,
		TemperatureHigh: make(map[string]float32),
		TemperatureLow:  make(map[string]float32),
		Condition:       make(map[string]Condition),
	}

	for i, key := range o.Daily.Time {
		w.TemperatureHigh[key] = o.Daily.Temperature2mMax[i]
		w.TemperatureLow[key] = o.Daily.Temperature2mMin[i]
		w.Condition[key] = Condition(o.Daily.WeatherCode[i])
	}

	return &w
}
