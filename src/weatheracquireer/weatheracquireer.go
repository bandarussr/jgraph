package weatheracquireer

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

func New(long, lat string, f Forecast) *Weather {
	w := &Weather{}

	m := getWeatherFromMetNo(lat, long)
	fmt.Printf("%f\n", m.Properties.Timeseries[0].Data.Instant.Details.AirTemperature)
	
	return w
}
