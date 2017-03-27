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

func TestAddPackageNameAndVersionToSignaturesFunction(t *testing.T) {
	a := CommitSignature{}
	a.Hash = "a"
	a.Email = "a@exapmle.com"
	a.Name = "Name A"
	a.Package = ""
	a.Signature = signature.PackageSignatures{
		"packageA": signature.Signature{
			Functions: []string{"func A() string"},
		},
	}
	a.Subject = "Subject A"
	a.Version = Version{1, 0, 0}

	b := CommitSignature{}
	b.Hash = "b"
	b.Email = "b@exapmle.com"
	b.Name = "Name B"
	b.Package = ""
	b.Signature = signature.PackageSignatures{
		"packageA": signature.Signature{
			Functions: []string{"func A() string"},
		},
	}
	b.Subject = "Subject A"
	b.Version = Version{0, 0, 1}

	signatures := []*CommitSignature{&a, &b}

	expectedPackageName := "github.com/a-h/example"
	addPackageNameAndVersionToSignatures(signatures, expectedPackageName)

	expectedVersion := Version{0, 0, 0}
	if a.Version != expectedVersion {
		t.Errorf("expected the first commit to have a version of 0.0.0, but was %v", a.Version)
	}
	if a.Package != expectedPackageName {
		t.Errorf("(1) expected package name %v, but was '%v'", expectedPackageName, a.Package)
	}

	expectedVersion = Version{0, 0, 1}
	if b.Version != expectedVersion {
		t.Errorf("expected the second commit to have a version of %v, but was %v", expectedVersion, b.Version)
	}
	if b.Package != expectedPackageName {
		t.Errorf("(2) expected package name %v, but was '%v'", expectedPackageName, b.Package)
	}
}
