package main

import (
	"fmt"
	"os"

	"github.com/bandarussr/jgraph/src/weather"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <latitude> <longitude>\n", os.Args[0]);
		return;
	}
	
	// Get weather information.
	w := weather.New(os.Args[1], os.Args[2], weather.ForecastDaily)
	fmt.Printf("%+v\n", w)
	// Create jgraph view.
}
