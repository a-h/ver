package main

import (
	"flag"
	"fmt"

	"github.com/a-h/ver/signature"

	"os"
)

var dir = flag.String("d", "", "The directory to analyse.")

func main() {
	flag.Parse()

	if *dir == "" {
		fmt.Println("Please provide a directory with the -d parameter.")
		os.Exit(-1)
	}

	signatures, err := signature.GetFromDirectory(*dir)

	if err != nil {
		fmt.Printf("Failed to get signatures of packages with error: %s\n", err.Error())
		os.Exit(-1)
	}

	for packageName, signature := range signatures {
		fmt.Printf("Package: %s\n", packageName)
		fmt.Println(signature)
	}
}
