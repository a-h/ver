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

	history, err := gitRepo.Log()

	if err != nil {
		fmt.Printf("Failed to get the git history with error: %s\n", err.Error())
		os.Exit(-1)
	}

	signatures := make([]CommitSignature, len(history))

	for idx, h := range history {
		fmt.Printf("Processing git log entry: %v\n", h)

		err := gitRepo.Get(h.Commit)

		if err != nil {
			fmt.Printf("Failed to get commit %s with error: %s\n", h.Commit, err.Error())
			break
		}

		cs, err := signature.GetFromDirectory(gitRepo.Location)

		if err != nil {
			fmt.Printf("Failed to get signatures of package at commit %s with error: %s\n", h.Commit, err.Error())
			os.Exit(-1)
		}

		signatures[idx] = CommitSignature{
			Commit:    h,
			Signature: cs,
		}
	}

	for _, cs := range signatures {
		fmt.Printf("Commit %s\n", cs.Commit.Commit)
		fmt.Println(cs.Signature)
		fmt.Println()
	}
}

type CommitSignature struct {
	Commit    history.History
	Signature signature.PackageSignatures
}
