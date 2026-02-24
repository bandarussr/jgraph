package weather

import "fmt"

const dev = true

type Weather struct {
	Location        string
	Keys            []string
	TemperatureHigh map[string]int
	TemperatureLow  map[string]int
	Condition       map[string]Condition
}

type Condition int

const (
	ConditionSun Condition = iota
	ConditionCloud
	ConditionRain
	ConditionSnow
	ConditionThunderstorm
	ConditionClear
)

type Forecast int

const (
	ForecastDaily Forecast = iota
	ForecastHourly
)

func New(lat, long string, f Forecast) *Weather {
	w := &Weather{}

	m := getWeatherFromOpenMeteo(lat, long)
	fmt.Printf("%f\n", m.Daily.Temperature2mMax[0])
	
	return w
}
