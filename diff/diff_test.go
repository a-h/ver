package diff

import (
	"reflect"
	"testing"

	"github.com/a-h/ver/signature"
)

func TestThatAMapCanBeCreatedFromAnArray(t *testing.T) {
	in := []string{"a", "b"}

	m := makeStringMap(in)

	if len(m) != 2 {
		t.Errorf("Expected the output map to have two keys, but it had %d keys", len(m))
	}

	if _, contains := m["a"]; !contains {
		t.Error("Expected the output map to contain 'a', but it was not found.")
	}

	if _, contains := m["b"]; !contains {
		t.Error("Expected the output map to contain 'b', but it was not found.")
	}
}

func TestThatStringArraysCanBeDiffed(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name     string
		current  []string
		next     []string
		expected Diff
	}{
		{
			name:    "Value removed",
			current: []string{"a"},
			next:    []string{},
			expected: Diff{
				Removed: 1,
			},
		},
		{
			name:    "Value added",
			current: []string{"a"},
			next:    []string{"b"},
			expected: Diff{
				Removed: 1,
				Added:   1,
			},
		},
		{
			name:    "No changes",
			current: []string{"a"},
			next:    []string{"a"},
			expected: Diff{
				Removed: 0,
				Added:   0,
			},
		},
		{
			name:    "Value added and removed",
			current: []string{"a"},
			next:    []string{"b"},
			expected: Diff{
				Removed: 1,
				Added:   1,
			},
		},
	}
	for _, tt := range tests {
		if actual := calculateStringDiff(tt.current, tt.next); !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("%q. Expected %v but got %v", tt.name, tt.expected, actual)
		}
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		current  signature.PackageSignatures
		next     signature.PackageSignatures
		expected SummaryDiff
	}{
		{
			name: "Function added to existing package",
			current: signature.PackageSignatures{
				"packageA": signature.Signature{
					Functions: []string{"func a() string"},
				},
			},
			next: signature.PackageSignatures{
				"packageA": signature.Signature{
					Functions: []string{"func a() string", "func b() string"},
				},
			},
			expected: SummaryDiff{
				PackageChanges: Diff{Added: 0, Removed: 0},
				Packages: []PackageDiff{
					PackageDiff{
						PackageName: "packageA",
						Functions:   Diff{Added: 1, Removed: 0},
					},
				},
			},
		},
		{
			name: "Constant added to existing package",
			current: signature.PackageSignatures{
				"packageA": signature.Signature{
					Functions: []string{"func a() string"},
				},
			},
			next: signature.PackageSignatures{
				"packageA": signature.Signature{
					Functions: []string{"func a() string"},
					Constants: []string{"const x = 0"},
				},
			},
			expected: SummaryDiff{
				PackageChanges: Diff{Added: 0, Removed: 0},
				Packages: []PackageDiff{
					PackageDiff{
						PackageName: "packageA",
						Functions:   Diff{Added: 0, Removed: 0},
						Constants:   Diff{Added: 1, Removed: 0},
					},
				},
			},
		},
		{
			name: "Package added and package removed",
			current: signature.PackageSignatures{
				"packageA": signature.Signature{
					Functions: []string{"func a() string"},
				},
			},
			next: signature.PackageSignatures{
				"packageB": signature.Signature{
					Functions: []string{"func b() string"},
				},
			},
			expected: SummaryDiff{
				PackageChanges: Diff{Added: 1, Removed: 1},
				Packages: []PackageDiff{
					PackageDiff{
						PackageName: "packageA",
						Functions:   Diff{Removed: 1},
					},
					PackageDiff{
						PackageName: "packageB",
						Functions:   Diff{Added: 1},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		actual := Calculate(tt.current, tt.next)

		if tt.expected.PackageChanges.Added != actual.PackageChanges.Added {
			t.Errorf("%q. Expected %d added packages, but calculated %d were added", tt.name, tt.expected.PackageChanges.Added, actual.PackageChanges.Added)
		}

		if tt.expected.PackageChanges.Removed != actual.PackageChanges.Removed {
			t.Errorf("%q. Expected %d removed packages, but calculated %d were removed", tt.name, tt.expected.PackageChanges.Removed, actual.PackageChanges.Removed)
		}

		if len(tt.expected.Packages) != len(actual.Packages) {
			t.Errorf("%q. Expected %d packages to be analysed, but got %d", tt.name, len(tt.expected.Packages), len(actual.Packages))
		}

		maxPackageCount := max(len(tt.expected.Packages), len(actual.Packages))
		for i := 0; i < maxPackageCount; i++ {
			exp, expPresent := getValueIfPossible(tt.expected.Packages, i)
			act, actPresent := getValueIfPossible(actual.Packages, i)

			if !expPresent {
				t.Errorf("%q. Couldn't compare package index %d because the expected value was missing.", tt.name, i)
				continue
			}

			if !actPresent {
				t.Errorf("%q. Couldn't compare package index %d because the actual value was missing.", tt.name, i)
				continue
			}

			if exp.PackageName != act.PackageName {
				t.Errorf("%q. For package index %d, expected package name of %s, but got %s", tt.name, i, exp.PackageName, act.PackageName)
			}

			testAreEqual(tt.name, i, "Constants", act.Constants, exp.Constants, t)
			testAreEqual(tt.name, i, "Fields", act.Fields, exp.Fields, t)
			testAreEqual(tt.name, i, "Functions", act.Functions, exp.Functions, t)
			testAreEqual(tt.name, i, "Interfaces", act.Interfaces, exp.Interfaces, t)
			testAreEqual(tt.name, i, "Structs", act.Structs, exp.Structs, t)
		}
	}
}

func testAreEqual(testName string, pkgIndex int, field string, actual Diff, expected Diff, t *testing.T) {
	if expected.Added != actual.Added {
		t.Errorf("%q. Package index %d: Expected %d %s added, but %d were found to have been added", testName, pkgIndex, expected.Added, field, actual.Added)
	}

	if expected.Removed != actual.Removed {
		t.Errorf("%q. Package index %d: Expected %d %s removed, but %d were found to have been removed", testName, pkgIndex, expected.Removed, field, actual.Removed)
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func getValueIfPossible(packageDiff []PackageDiff, index int) (diff PackageDiff, ok bool) {
	if index >= len(packageDiff) {
		return diff, false
	}

	return packageDiff[index], true
}
