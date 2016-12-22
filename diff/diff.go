package diff

import "github.com/a-h/ver/signature"

// SummaryDiff provides a summary of changes to a set of packages.
type SummaryDiff struct {
	PackageChanges Diff          `json:"packageChanges"`
	Packages       []PackageDiff `json:"packages"`
}

// PackageDiff describes changes to a given package.
type PackageDiff struct {
	PackageName string `json:"packageName"`
	Functions   Diff   `json:"functions"`
	Fields      Diff   `json:"fields"`
	Constants   Diff   `json:"constants"`
	Structs     Diff   `json:"structs"`
	Interfaces  Diff   `json:"interfaces"`
}

// Diff describes the changes to an element (added, removed).
type Diff struct {
	Removed int `json:"removed"`
	Added   int `json:"added"`
}

// Calculate the difference between package signatures.
func Calculate(current signature.PackageSignatures, next signature.PackageSignatures) SummaryDiff {
	d := &SummaryDiff{}
	for currPkgKey, currPkgSig := range current {
		nextPkgSig, ok := next[currPkgKey]

		if !ok {
			// Package is missing, if it is missing, calculating the package diff
			// is based on the zero value of a PackageSignature
			d.PackageChanges.Removed++
		}

		d.Packages = append(d.Packages, PackageDiff{
			PackageName: currPkgKey,
			Constants:   calculateStringDiff(currPkgSig.Constants, nextPkgSig.Constants),
			Fields:      calculateStringDiff(currPkgSig.Fields, nextPkgSig.Fields),
			Functions:   calculateStringDiff(currPkgSig.Functions, nextPkgSig.Functions),
			Interfaces:  calculateStringDiff(currPkgSig.Interfaces, nextPkgSig.Interfaces),
			Structs:     calculateStringDiff(currPkgSig.Structs, nextPkgSig.Structs),
		})
	}

	for nextPkgKey, nextPkgSig := range next {
		_, ok := current[nextPkgKey]

		if ok {
			// We've already compared this package, since it exists in the current version.
			continue
		}

		// We have a new package.
		d.PackageChanges.Added++

		// Since we have a completely new package, everything is new.
		d.Packages = append(d.Packages, PackageDiff{
			PackageName: nextPkgKey,
			Constants:   Diff{Added: len(nextPkgSig.Constants)},
			Fields:      Diff{Added: len(nextPkgSig.Fields)},
			Functions:   Diff{Added: len(nextPkgSig.Functions)},
			Interfaces:  Diff{Added: len(nextPkgSig.Interfaces)},
			Structs:     Diff{Added: len(nextPkgSig.Structs)},
		})
	}

	return *d
}

func calculateStringDiff(current []string, next []string) Diff {
	c := makeStringMap(current)
	n := makeStringMap(next)

	d := &Diff{}

	for currentKey := range c {
		if _, ok := n[currentKey]; !ok {
			d.Removed++
		}
	}

	for nextKey := range n {
		if _, ok := c[nextKey]; !ok {
			d.Added++
		}
	}

	return *d
}

func makeStringMap(a []string) map[string]bool {
	m := make(map[string]bool, len(a))

	for _, v := range a {
		m[v] = true
	}

	return m
}
