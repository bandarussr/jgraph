package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type metNoResponse struct {
	Properties struct {
		Timeseries []struct {
			Time string `json:"time"`
			Data struct {
				Instant struct {
					Details struct {
						AirTemperature   float64 `json:"air_temperature"`
						RelativeHumidity float64 `json:"relative_humidity"`
						WindSpeed        float64 `json:"wind_speed"`
					} `json:"details"`
				} `json:"instant"`
				Next12Hours *struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
					} `json:"details"`
				} `json:"next_12_hours"`
				Next1Hours *struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
					} `json:"details"`
				} `json:"next_1_hours"`
				Next6Hours *struct {
					Details struct {
						TempMax float64 `json:"air_temperature_max"`
						TempMin float64 `json:"air_temperature_min"`
					} `json:"details"`
				} `json:"next_6_hours"`
			} `json:"data"`
		} `json:"timeseries"`
	} `json:"properties"`
}

func getWeatherFromMetNo(long, lat string) *metNoResponse {
	// Limit API requests during development.
	if dev {
		data, err := os.ReadFile(".devdata.json")
		if err == nil {
			var weather metNoResponse
			err = json.Unmarshal(data, &weather)
			if err == nil {
				return &weather
			}
		}
	}
	
	url := fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%s&lon=%s", lat, long)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Error creating request for weather data.")
	}
	
	req.Header.Set("User-Agent", "JgraphLab/1.0 github.com/bandarussr/jgraph")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic("Error getting weather data.")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic("Error getting weather data: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic("Error reading weather data.")
	}
	
	var weather metNoResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic("Error pasing weather data: " + err.Error())
	}
	
	if dev {
		err = os.WriteFile(".devdata.json", body, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing weather data to file: %s\n", err.Error())
		}
	}
	
	return &weather
}
