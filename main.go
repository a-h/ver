package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/a-h/ver/diff"
	"github.com/a-h/ver/git"
	"github.com/a-h/ver/signature"

	"encoding/json"
	"net/url"
	"os"
)

var repo = flag.String("r", "", "The git repo to clone and analyse, e.g. https://github.com/a-h/ver")
var out = flag.String("o", "", "When set, outputs to a file in JSON format.")

func main() {
	flag.Parse()

	if *repo == "" {
		fmt.Println("Please provide a repo with the -r parameter.")
		os.Exit(-1)
	}

	repoURL, err := url.Parse(strings.ToLower(*repo))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse URL %v: %v\n", *repo, err)
		os.Exit(-1)
	}

	var outFile *os.File
	if *out != "" {
		outFile, err = os.Create(*out)
		defer outFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open output file: %v\n", err)
			os.Exit(-1)
		}
	}

	gitRepo, err := git.Clone(*repo)
	defer gitRepo.CleanUp()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to clone git repo: %v\n", err)
		os.Exit(-1)
	}

	if err = gitRepo.Fetch(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch from git repo: %v\n", err)
		os.Exit(-1)
	}

	fmt.Printf("Cloned repo %s into %s\n", *repo, gitRepo.PackageDirectory())

	history, err := gitRepo.Log()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get the git history: %v\n", err)
		os.Exit(-1)
	}

	signatures := make([]*CommitSignature, len(history))

	fatalError := false

	for idx, h := range history {
		fmt.Printf("Processing git log entry: %v\n", h)

		cs := &CommitSignature{
			Commit: h,
		}

		err := gitRepo.Get(h.Hash)

		if err != nil {
			cs.Error = fmt.Errorf("Failed to get commit %s: %s\n", h.Hash, err.Error())
			signatures[idx] = cs
			continue
		}

		err = goget(gitRepo.BaseLocation, gitRepo.PackageDirectory())

		if err != nil {
			cs.Error = err
			signatures[idx] = cs
			continue
		}

		sig, err := signature.GetFromDirectory(gitRepo.BaseLocation, gitRepo.PackageDirectory())

		if err != nil {
			cs.Error = fmt.Errorf("Failed to get signatures of package at commit %s: %s\n",
				h.Hash, err.Error())
			signatures[idx] = cs
			continue
		}

		cs.Signature = sig
		signatures[idx] = cs

		err = gitRepo.Revert()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to revert the repo back to HEAD: %s\n", err.Error())
			fatalError = true
			break
		}
	}

	if fatalError {
		fmt.Fprintf(os.Stderr, "Failed with fatal error.\n")
		os.Exit(-1)
		return
	}

	fmt.Printf("About to calculate signatures...\n")

	addPackageNameAndVersionToSignatures(signatures, repoURL.Host+repoURL.Path)

	for _, cs := range signatures {
		fmt.Println()
		fmt.Printf("Commit: %s\n", cs.Commit.Hash)
		fmt.Printf("Subject: %s\n", cs.Commit.Subject)
		fmt.Printf("Date: %v\n", cs.Commit.Date())
		fmt.Printf("Version: %v\n", cs.Version)
		if err != nil {
			fmt.Printf("Error: %v\n", cs.Error)
		}
		if outFile != nil {
			j, err := json.Marshal(cs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to marshal JSON output: %v", err)
			}
			_, err = outFile.Write(j)
			outFile.Write([]byte{0x0A})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write to output file: %v", err)
			}
		}
	}
}

func goget(gopath string, location string) error {
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
		return fmt.Errorf("failed to go get all in repo with gopath %s at directory %s: '%v' message '%s'", gopath, location, err, string(out))
	}

	return nil
}

func addPackageNameAndVersionToSignatures(signatures []*CommitSignature, packageName string) {
	version := Version{}

	if len(signatures) > 0 {
		previous := signatures[0]
		previous.Package = packageName
		previous.Version = version

		for _, current := range signatures[1:] {
			current.Package = packageName
			if current.Error != nil {
				// Add 1 to the build, even though it wasn't successfully handled.
				version = version.Add(Version{Build: 1})
				current.Version = version
				continue
			}

			// Calculate the diff against the previous version.
			diff := diff.Calculate(previous.Signature, current.Signature)
			// Work out what the version increment should be.
			delta := calculateVersionDelta(diff)
			version = version.Add(delta)
			current.Version = version

			// Update the previous version.
			previous = current
		}
	}
}

func calculateVersionDelta(sd diff.SummaryDiff) Version {
	d := Version{
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

// CommitSignature is the signature of a commit.
type CommitSignature struct {
	git.Commit
	Package   string                      `json:"pkg"`
	Signature signature.PackageSignatures `json:"-"`
	Error     error                       `json:"error"`
	Version   Version                     `json:"v"`
}
