package main

import (
	"flag"
	"fmt"

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

	for _, cs := range signatures {
		if cs.Error != nil {
			continue
		}

		fmt.Println()
		fmt.Printf("Commit %s\n", cs.Commit.Commit)

		for k, v := range cs.Signature {
			fmt.Printf("Package %s { Constants: %d, Fields: %d, Functions: %d, Interfaces: %d, Structs %d }\n", k,
				len(v.Constants), len(v.Fields), len(v.Functions), len(v.Interfaces), len(v.Structs))
		}
		fmt.Println()
	}
}

// Version represents a major, minor and build version.
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Build int `json:"build"`
}

// CommitSignature is the signature of a commit.
type CommitSignature struct {
	Commit    history.History             `json:"commit"`
	Signature signature.PackageSignatures `json:"signature"`
	Error     error                       `json:"error"`
}
