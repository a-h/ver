package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/a-h/ver/diff"
	"github.com/a-h/ver/git"
	"github.com/a-h/ver/measure"
	"github.com/a-h/ver/signature"

	"os"
)

var repo = flag.String("r", "", "The git repo to clone and analyse, e.g. https://github.com/a-h/ver")

func main() {
	flag.Parse()

	if *repo == "" {
		fmt.Println("Please provide a repo with the -r parameter.")
		os.Exit(-1)
	}

	gitRepo, err := git.Clone(*repo)
	defer gitRepo.CleanUp()

	if err != nil {
		fmt.Printf("Failed to clone git repo with error: %s\n", err.Error())
		os.Exit(-1)
	}

	if err = gitRepo.Fetch(); err != nil {
		fmt.Printf("Failed to fetch from git repo with error: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Printf("Cloned repo %s into %s\n", *repo, gitRepo.PackageDirectory())

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

		err := gitRepo.Get(h.Hash)

		if err != nil {
			cs.Error = fmt.Errorf("Failed to get commit %s with error: %s\n", h.Hash, err.Error())
			signatures[idx] = *cs
			continue
		}

		err = goget(gitRepo.BaseLocation, gitRepo.PackageDirectory())

		if err != nil {
			cs.Error = err
			signatures[idx] = *cs
			continue
		}

		sig, err := signature.GetFromDirectory(gitRepo.BaseLocation, gitRepo.PackageDirectory())

		if err != nil {
			cs.Error = fmt.Errorf("Failed to get signatures of package at commit %s with error: %s\n",
				h.Hash, err.Error())
			signatures[idx] = *cs
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
		fmt.Printf("Failed with fatal error.\n")
		return
	}

	fmt.Printf("About to calcaulate signatures...\n")

	calculateVersionsFromSignatures(signatures)
}

func goget(gopath string, location string) error {
	defer measure.TimeTrack(time.Now(), "main.goget")

	os.Chdir(location)

	// Set the path, the Go tools take the first GOPATH in the set.
	env := os.Environ()

	for i, v := range env {
		if strings.HasPrefix(v, "GOPATH=") {
			env[i] = fmt.Sprintf("GOPATH=%s", gopath)
			break
		}
	}

	// Get the stuff.
	cmd := exec.Command("go", "get", "./...")
	cmd.Env = env
	cmd.Dir = location
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to go get all in repo with gopath %s at directory %s with err '%v' message '%s'", gopath, location, err, string(out))
	}

	return nil
}

func calculateVersionsFromSignatures(signatures []CommitSignature) {
	defer measure.TimeTrack(time.Now(), "main.calculateVersionsFromSignatures")

	version := Version{}

	if len(signatures) > 0 {
		current := signatures[0]

		for _, cs := range signatures[1:] {
			if cs.Error != nil {
				// Add 1 to the build, even though it wasn't successfully handled.
				version = addDeltaToVersion(version, Version{Build: 1})

				fmt.Println()
				fmt.Printf("Commit: %s\n", cs.Commit.Hash)
				fmt.Printf("Commit: %s\n", cs.Commit.Subject)
				fmt.Printf("Date: %s\n", cs.Commit.Date())
				fmt.Printf("Error: %s\n", cs.Error.Error())
				continue
			}

			// Calculate the diff against the previous version.
			diff := diff.Calculate(current.Signature, cs.Signature)
			// Work out what the version increment should be.
			delta := calculateVersionDelta(diff)
			version = addDeltaToVersion(version, delta)

			fmt.Println()
			fmt.Printf("Commit: %s\n", cs.Commit.Hash)
			fmt.Printf("Subject: %s\n", cs.Commit.Subject)
			fmt.Printf("Date: %s\n", cs.Commit.Date())
			fmt.Printf("Version: %s\n", version.String())
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

func calculateVersionDelta(sd diff.SummaryDiff) Version {
	defer measure.TimeTrack(time.Now(), "main.calculateVersionDelta")

	d := &Version{
		Build: 1, // Always increment the build.
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

	return *d
}

func updateBasedOn(d diff.Diff, binaryCompatibilityBroken *bool, newExportedData *bool) {
	if d.Added > 0 {
		*newExportedData = true
	}

	if d.Removed > 0 {
		*binaryCompatibilityBroken = true
	}
}

// CommitSignature is the signature of a commit.
type CommitSignature struct {
	Commit    git.Commit                  `json:"commit"`
	Signature signature.PackageSignatures `json:"signature"`
	Error     error                       `json:"error"`
}
