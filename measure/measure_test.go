package measure

import (
	"testing"
	"time"
)

func TestTimeTrack(t *testing.T) {	
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	elapsedTime := TimeTrack(start, "TestTimeTrack")

	if !(elapsedTime >= 100*time.Millisecond) {
		t.Errorf("ElapsedTime should be greater than 100ms")
	}
}
