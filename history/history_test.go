package history

import (
	"reflect"
	"testing"
	"time"
)

func TestThatARepoCanBeCloned(t *testing.T) {
	if testing.Short() {
		t.Skip("Cloning a repo requires network access.")
	}

	h, err := Clone("https://github.com/a-h/ver")

	if err != nil {
		t.Fatal(err)
	}

	defer h.CleanUp()

	log, err := h.Log()

	if err != nil {
		t.Fatal(err)
	}

	expected := []History{
		{
			Commit:  "f5ea0f3b4f65fa179967d4d4d4709662ffc711b8",
			Subject: "First-commit",
			Name:    "Adrian Hesketh",
			Email:   "adrianhesketh@hushmail.com",
			Date:    time.Date(2016, 12, 17, 15, 47, 57, 0, time.Local),
		},
	}

	for idx, a := range log {
		if idx < len(expected) {
			e := expected[idx]

			if !reflect.DeepEqual(e, a) {
				t.Errorf("Expected a commit of %v, but got %v", e, a)
			}
		}
	}
}
