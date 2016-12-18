package history

import (
	"fmt"
	"testing"
)

func TestThatARepoCanBeCloned(t *testing.T) {
	h, err := Clone("https://github.com/a-h/ver")

	if err != nil {
		t.Fatal(err)
	}

	defer h.CleanUp()

	log, err := h.Log()

	if err != nil {
		t.Fatal(err)
	}

	for _, l := range log {
		fmt.Println(l)
	}
}
