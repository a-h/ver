package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var in = flag.String("in", "github.csv", "The name of the CSV file containing repository names.")

func main() {
	flag.Parse()

	f, err := os.Open(*in)
	if err != nil {
		fmt.Println("Failed!")
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println("github.com/" + scanner.Text())
	}
}
