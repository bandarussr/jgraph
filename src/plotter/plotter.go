package plotter

import (
	_ "embed"
	"os"
	"text/template"

	"github.com/bandarussr/jgraph/src/weather"
)

//go:embed template.jgr
var jgr_template string

const XMin float32 = 20
const XMax float32 = 75
const YMin float32 = -0.5
const YMax float32 = 9.5

type JGraph struct {
	XMin float32
	XMax float32
	YMin float32
	YMax float32

	LowCol  float32
	HighCol float32

	Location     string
	ForecastType string
	Forecast     []Forecast
}

type Forecast struct {
	Key     string
	Row     float32
	Low     int
	High    int
	TempBar []TempBar
}

type TempBar struct {
	R     float32
	G     float32
	B     float32
	Start float32
	End   float32
}

func New(w weather.Weather) *JGraph {
	j := &JGraph{
		XMin: XMin,
		XMax: XMax,
		YMin: YMin,
		YMax: YMax,

		LowCol:  27,
		HighCol: 69,

		Location:     w.Location,
		ForecastType: "Daily",
		Forecast:     make([]Forecast, len(w.Keys)),
	}

	for i, key := range w.Keys {
		f := &j.Forecast[i]
		f.Key = key
		f.Row = float32(len(w.Keys) - 1 - i)
		f.Low = int(w.TemperatureLow[key])
		f.High = int(w.TemperatureHigh[key])
		f.TempBar = f.makeTempBar(j.LowCol, j.HighCol, int(w.MinTemperature), int(w.MaxTemperature))
	}

	return j
}

func (j *JGraph) Plot(file string) {
	tmpl, err := template.New("jgraph").Parse(jgr_template)
	if err != nil {
		panic("Error parsing template: " + err.Error())
	}

	f, err := os.Create(file)
	if err != nil {
		panic("Error creating output file: " + err.Error())
	}

	tmpl.Execute(f, j)
}

// ALL BELOW CODED BY CLAUDE ..CHECK LATER!
func (f *Forecast) makeTempBar(lowCol, highCol float32, minTemp, maxTemp int) []TempBar {
	steps := 80
	bars := make([]TempBar, steps)
	stepSize := (highCol - lowCol) / float32(steps)

	tempRange := float32(maxTemp - minTemp)
	xRange := highCol - lowCol

	for i := range steps {
		x1 := lowCol + float32(i)*stepSize
		x2 := x1 + stepSize + 0.05

		xNorm := (x1 - lowCol) / xRange
		temp := float32(minTemp) + xNorm*tempRange

		t := (temp - float32(minTemp)) / tempRange
		t = clamp(t, 0.0, 1.0)

		var r, g, b float32
		if temp < float32(f.Low) || temp > float32(f.High) {
			// outside this day's range — gray
			r, g, b = 0.4, 0.4, 0.4
		} else {
			r, g, b = tempToColor(t)
		}

		bars[i] = TempBar{
			R:     r,
			G:     g,
			B:     b,
			Start: x1,
			End:   x2,
		}
	}

	return bars
}

// Magic color gradient, no clue how this works.
func tempToColor(t float32) (r, g, b float32) {
	stops := [][4]float32{
		{0.0, 0.3, 0.7, 0.95},   // blue
		{0.33, 0.2, 0.85, 0.85}, // cyan
		{0.66, 0.4, 0.85, 0.5},  // green
		{1.0, 0.75, 0.85, 0.2},  // yellow-green
	}

	for i := range len(stops) - 1 {
		t0, c0 := stops[i][0], stops[i][1:]
		t1, c1 := stops[i+1][0], stops[i+1][1:]
		if t <= t1 {
			f := (t - t0) / (t1 - t0)
			return c0[0] + f*(c1[0]-c0[0]),
				c0[1] + f*(c1[1]-c0[1]),
				c0[2] + f*(c1[2]-c0[2])
		}
	}

	last := stops[len(stops)-1]
	return last[1], last[2], last[3]
}

func clamp(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
