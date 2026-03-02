package plotter

import (
	_ "embed"
	"os"
	"text/template"

	"github.com/bandarussr/jgraph/src/weather"
)

//go:embed template.jgr
var jgr_template string

const XMin float32 = 0
const XMax float32 = 100
const YMin float32 = -1

type JGraph struct {
	XMin float32
	XMax float32
	YMin float32
	YMax float32

	TitleY       float32
	DayCol       float32
	ConditionCol float32
	LowCol       float32
	HighCol      float32
	TempBarMin   float32
	TempBarMax   float32

	Location string
	Forecast []Forecast
}

type Forecast struct {
	Key       string
	Row       float32
	Low       int
	High      int
	TempBar   []TempBar
	Condition string
}

type TempBar struct {
	R     float32
	G     float32
	B     float32
	Start float32
	End   float32
}

func New(w weather.Weather) *JGraph {
	tempBarBuffer := float32(3)

	j := &JGraph{
		XMin: XMin,
		XMax: XMax,
		YMin: YMin,
		YMax: float32(len(w.Keys)),

		TitleY:       7,
		DayCol:       -25,
		ConditionCol: -15,
		LowCol:       -1,
		HighCol:      101,
		TempBarMin:   XMin + tempBarBuffer,
		TempBarMax:   XMax - tempBarBuffer,

		Location: w.Location,
		Forecast: make([]Forecast, len(w.Keys)),
	}

	// Condition to symbol file
	conditionSymbol := map[weather.Condition]string{
		weather.ConditionSun:          "symbols/condition_sun.eps",
		weather.ConditionCloud:        "symbols/condition_cloud.eps",
		weather.ConditionRain:         "symbols/condition_rain.eps",
		weather.ConditionSnow:         "symbols/condition_snow.eps",
		weather.ConditionThunderstorm: "symbols/condition_thunderstorm.eps",
	}

	for i, key := range w.Keys {
		f := &j.Forecast[i]
		f.Key = key
		f.Row = float32(len(w.Keys) - 1 - i)
		f.Low = int(w.TemperatureLow[key])
		f.High = int(w.TemperatureHigh[key])
		f.TempBar = f.makeTempBar(j.LowCol+tempBarBuffer, j.HighCol-tempBarBuffer, int(w.MinTemperature), int(w.MaxTemperature))
		f.Condition = conditionSymbol[w.Condition[key]]
	}

	return j
}

func (j *JGraph) Plot() {
	tmpl, err := template.New("jgraph").Parse(jgr_template)
	if err != nil {
		panic("Error parsing template: " + err.Error())
	}

	err = tmpl.Execute(os.Stdout, j)
	if err != nil {
	    panic("Error executing template: " + err.Error())
	}
}

func (f *Forecast) makeTempBar(lowCol, highCol float32, globalMin, globalMax int) []TempBar {
	steps := 80
	bars := make([]TempBar, steps)

	tempRange := float32(globalMax - globalMin)
	xRange := highCol - lowCol

	// X positions where this day's colored bar starts and ends
	dayLowX := lowCol + (float32(f.Low-globalMin)/tempRange)*xRange
	dayHighX := lowCol + (float32(f.High-globalMin)/tempRange)*xRange

	stepSize := (dayHighX - dayLowX) / float32(steps)

	for i := range steps {
		x1 := dayLowX + float32(i)*stepSize
		x2 := x1 + stepSize + 0.05

		// t is 0.0 at global min, 1.0 at global max
		t := (x1 - lowCol) / xRange
		if t < 0.0 {
			t = 0.0
		} else if t > 1.0 {
			t = 1.0
		}

		r, g, b := tempToColor(t)
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
