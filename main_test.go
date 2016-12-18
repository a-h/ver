package main

import (
	"path"
	"testing"

	"strings"

	"github.com/a-h/ver/signature"
)
import "os"

func TestThatItemsAreExtracted(t *testing.T) {
	wd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Couldn't find path to example with error %v", err)
	}

	ps, err := signature.GetFromDirectory(path.Join(wd, "example"))

	if err != nil {
		t.Fatal("failed to get subpackages: " + err.Error())
	}

	if len(ps) != 2 {
		t.Errorf("Expected %d packages to be discovered, got %d packages", 2, len(ps))
	}

	for name, signature := range ps {
		if strings.HasSuffix(name, "example") {
			assert("example Constants", 1, len(signature.Constants), t)
			assert("example Fields", 1, len(signature.Fields), t)
			assert("example Functions", 6, len(signature.Functions), t)
			assert("example Interfaces", 1, len(signature.Interfaces), t)
			assert("example Structs", 4, len(signature.Structs), t)
			continue
		}

		if strings.HasSuffix(name, "example2") {
			assert("example2 Constants", 0, len(signature.Constants), t)
			assert("example2 Fields", 1, len(signature.Fields), t)
			assert("example2 Functions", 0, len(signature.Functions), t)
			assert("example2 Interfaces", 0, len(signature.Interfaces), t)
			assert("example2 Structs", 1, len(signature.Structs), t)
			continue
		}

		t.Fatalf("Expected example and example2 - didn't expect the path %s", name)
	}
}

func assert(name string, expected int, actual int, t *testing.T) {
	if actual != expected {
		t.Errorf("Test '%s' failed, expected %d, got %d", name, expected, actual)
	}
}
