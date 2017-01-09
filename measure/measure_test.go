package measure

import (
	"testing"
	"time"
)

func TestTimeTrack(t *testing.T) {

	start := time.Now()

	time.Sleep(100 * time.Millisecond)

	elapsedTime := TimeTrack(start, "TestTimeTrack")

	if !(elapsedTime >= 99*time.Millisecond && elapsedTime <= 101*time.Millisecond) {
		t.Errorf("ElapsedTime should be between 99ms and 101 ms")
	}
}
