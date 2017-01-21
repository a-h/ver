package measure

import (
	"log"
	"os"
	"time"
)

// TimeTrack calculates the seconds elapsed for a named function,
// it's enabled using the `MEASURE_PERF` environment variable.
// For example, used as below at the start of a function:
// `defer measure.TimeTrack(time.Now(), "Get")``
func TimeTrack(start time.Time, name string) time.Duration {
	elapsed := time.Since(start)
	if os.Getenv("MEASURE_PERF") != "" {
		log.Printf("%s took %s", name, elapsed)
	}
	return elapsed
}
