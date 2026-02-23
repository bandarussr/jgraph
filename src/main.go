package main

import (
	"fmt"
	"os"

	"github.com/bandarussr/jgraph/src/weatheracquireer"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <longitude> <latitude>\n", os.Args[0]);
		return;
	}
	
	// Get weather information.
	weatheracquireer.New(os.Args[1], os.Args[2], weatheracquireer.ForecastDaily)
	
	// Create jgraph view.
}
