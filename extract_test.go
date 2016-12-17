package main

import "testing"

func TestThatItemsAreExtracted(t *testing.T) {
	pn, err := getSubpackages("github.com/a-h/ver/example")

	if err != nil {
		t.Fatal("failed to get subpackages")
	}

	if len(pn) != 2 {
		t.Errorf("Expected %d packages to be discovered, got %d packages", 2, len(pn))
	}
}
