package main

import (
	"flag"
	"fmt"

	"github.com/a-h/ver/diff"
	"github.com/a-h/ver/history"
	"github.com/a-h/ver/signature"

	"os"
)

var repo = flag.String("r", "", "The git repo to analyse.")

func main() {
	flag.Parse()

	if *repo == "" {
		fmt.Println("Please provide a repo with the -r parameter.")
		os.Exit(-1)
	}

	gitRepo, err := history.Clone(*repo)
	defer gitRepo.CleanUp()

	if err != nil {
		fmt.Printf("Failed to clone git repo with error: %s\n", err.Error())
		os.Exit(-1)
	}

	if err = gitRepo.Fetch(); err != nil {
		fmt.Printf("Failed to fetch from git repo with error: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Printf("Cloned repo %s into %s\n", *repo, gitRepo.Location)

	history, err := gitRepo.Log()

	if err != nil {
		fmt.Printf("Failed to get the git history with error: %s\n", err.Error())
		os.Exit(-1)
	}

	signatures := make([]CommitSignature, len(history))

	fatalError := false

	for idx, h := range history {
		fmt.Printf("Processing git log entry: %v\n", h)

		cs := &CommitSignature{
			Commit: h,
		}

		err := gitRepo.Get(h.Commit)

		if err != nil {
			cs.Error = fmt.Errorf("Failed to get commit %s with error: %s\n", h.Commit, err.Error())
			signatures[idx] = *cs
			continue
		}

		sig, err := signature.GetFromDirectory(gitRepo.Location)

		if err != nil {
			cs.Error = fmt.Errorf("Failed to get signatures of package at commit %s with error: %s\n", h.Commit, err.Error())
			continue
		}

		cs.Signature = sig
		signatures[idx] = *cs

		err = gitRepo.Revert()

		if err != nil {
			fmt.Printf("Failed to revert the repo back to HEAD with error: %s\n", err.Error())
			fatalError = true
			break
		}
	}

	if fatalError {
		return
	}

	version := Version{}

	if len(signatures) > 0 {
		current := signatures[0]

		for _, cs := range signatures[1:] {
			if cs.Error != nil {
				continue
			}

			// Calculate the diff and add it to the history.
			diff := diff.Calculate(current.Signature, cs.Signature)
			delta := calculateDelta(diff)
			version := addDeltaToVersion(version, delta)

			fmt.Println()
			fmt.Printf("Commit %s\n", cs.Commit.Commit)
			fmt.Printf("Version: %s\n", version.String())

			for k, v := range cs.Signature {
				fmt.Printf("Package %s { Constants: %d, Fields: %d, Functions: %d, Interfaces: %d, Structs %d }\n", k,
					len(v.Constants), len(v.Fields), len(v.Functions), len(v.Interfaces), len(v.Structs))
			}
			fmt.Println()
		}
	}
}

func addDeltaToVersion(v Version, d Version) Version {
	return Version{
		Major: v.Major + d.Major,
		Minor: v.Minor + d.Minor,
		Build: v.Build + d.Build,
	}
}

func calculateDelta(sd diff.SummaryDiff) Version {
	d := Version{
		Minor: 1, // Always increment the build.
	}

	binaryCompatibilityBroken := false
	newExportedData := false

	if sd.PackageChanges.Added > 0 {
		newExportedData = true
	}

	if sd.PackageChanges.Removed > 0 {
		binaryCompatibilityBroken = true
	}

	for _, pkg := range sd.Packages {
		updateBasedOn(pkg.Constants, &binaryCompatibilityBroken, &newExportedData)
		updateBasedOn(pkg.Fields, &binaryCompatibilityBroken, &newExportedData)
		updateBasedOn(pkg.Functions, &binaryCompatibilityBroken, &newExportedData)
		updateBasedOn(pkg.Interfaces, &binaryCompatibilityBroken, &newExportedData)
		updateBasedOn(pkg.Structs, &binaryCompatibilityBroken, &newExportedData)
	}

	if binaryCompatibilityBroken {
		d.Major = 1
	}

	if newExportedData {
		d.Minor = 1
	}

	return d
}

func updateBasedOn(d diff.Diff, binaryCompatibilityBroken *bool, newExportedData *bool) {
	if d.Added > 0 {
		*newExportedData = true
	}

	if d.Removed > 0 {
		*binaryCompatibilityBroken = true
	}
}

// Version represents a major, minor and build version.
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Build int `json:"build"`
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Build)
}

// CommitSignature is the signature of a commit.
type CommitSignature struct {
	Commit    history.History             `json:"commit"`
	Signature signature.PackageSignatures `json:"signature"`
	Error     error                       `json:"error"`
}
