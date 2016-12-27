package main

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/a-h/ver/diff"
	"github.com/a-h/ver/signature"
)

func TestThatItemsAreExtracted(t *testing.T) {
	wd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Couldn't find path to example with error %v", err)
	}

	ps, err := signature.GetFromDirectory(os.Getenv("GOPATH"), path.Join(wd, "example"))

	if err != nil {
		t.Fatal("failed to get subpackages: " + err.Error())
	}

	if len(ps) != 2 {
		t.Errorf("Expected 2 packages to be discovered, got %d packages", len(ps))
	}

	for name, signature := range ps {
		if strings.HasSuffix(name, "example") {
			assert("example Constants", 1, len(signature.Constants), t)
			assert("example Fields", 1, len(signature.Fields), t)
			assert("example Functions", 7, len(signature.Functions), t)
			assert("example Interfaces", 1, len(signature.Interfaces), t)
			assert("example Structs", 4, len(signature.Structs), t)
			continue
		}

		if strings.HasSuffix(name, "example2") {
			assert("example Constants", 0, len(signature.Constants), t)
			assert("example Fields", 1, len(signature.Fields), t)
			assert("example Functions", 0, len(signature.Functions), t)
			assert("example Interfaces", 0, len(signature.Interfaces), t)
			assert("example Structs", 1, len(signature.Structs), t)
			continue
		}

		t.Fatalf("Expected example - didn't expect the path %s", name)
	}
}

func assert(name string, expected int, actual int, t *testing.T) {
	if actual != expected {
		t.Errorf("Test '%s' failed, expected %d, got %d", name, expected, actual)
	}
}

func TestThatBinaryCompatibilityAndNewExportedDataCanBeUpdated(t *testing.T) {
	binaryCompatibilityBroken := false
	newExportedData := false
	updateBasedOn(diff.Diff{Added: 1, Removed: 1}, &binaryCompatibilityBroken, &newExportedData)

	if !binaryCompatibilityBroken {
		t.Errorf("Failed to update the binaryCompatibilityBroken variable.")
	}

	if !newExportedData {
		t.Errorf("Failed to update the newExportedData variable.")
	}
}

func TestThatVersionsCanBeUpdated(t *testing.T) {
	v2 := addDeltaToVersion(Version{Major: 1, Minor: 1, Build: 1}, Version{Major: 1, Minor: 1, Build: 1})

	if v2.String() != "2.2.2" {
		t.Errorf("Expected %s, but got %s", "2.2.2", v2.String())
	}
}

func TestThatVersionDeltasCanBeCalculated(t *testing.T) {
	tests := []struct {
		name     string
		sd       diff.SummaryDiff
		expected Version
	}{
		{
			name: "Package removed",
			sd: diff.SummaryDiff{
				PackageChanges: diff.Diff{Removed: 1},
			},
			expected: Version{Major: 1, Minor: 0, Build: 1},
		},
		{
			name: "Package added",
			sd: diff.SummaryDiff{
				PackageChanges: diff.Diff{Added: 1},
			},
			expected: Version{Major: 0, Minor: 1, Build: 1},
		},
		{
			name: "Function removed",
			sd: diff.SummaryDiff{
				Packages: []diff.PackageDiff{
					diff.PackageDiff{
						Functions: diff.Diff{Removed: 1},
					},
				},
			},
			expected: Version{Major: 1, Minor: 0, Build: 1},
		},
		{
			name: "Function added",
			sd: diff.SummaryDiff{
				Packages: []diff.PackageDiff{
					diff.PackageDiff{
						Functions: diff.Diff{Added: 1},
					},
				},
			},
			expected: Version{Major: 0, Minor: 1, Build: 1},
		},
	}
	for _, tt := range tests {
		if actual := calculateVersionDelta(tt.sd); tt.expected != actual {
			t.Errorf("%q. Expected version %s, but got %s", tt.name, tt.expected, actual)
		}
	}
}
