package main

import "testing"
import "fmt"

func TestThatItemsAreExtracted(t *testing.T) {
	pn, err := getSubpackages("github.com/a-h/ver/example")

	if err != nil {
		t.Fatal("failed to get subpackages")
	}

	if len(pn) != 2 {
		t.Errorf("Expected %d packages to be discovered, got %d packages", 2, len(pn))
	}

	info, err := getInformationFromPackages(pn)

	if err != nil {
		t.Fatal("failed to parse packages")
	}

	//TODO: Validate that all the public types are there.
	for _, item := range info {
		fmt.Println(item)
	}
}
