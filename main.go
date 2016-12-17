package main

import (
	"flag"
	"fmt"
	"go/types"
	"path/filepath"

	"strings"

	"os"

	"path"

	"bytes"

	"golang.org/x/tools/go/loader"
)

var packageName = flag.String("p", "github.com/a-h/ver/example", "The package to analyse, e.g. 'github.com/a-h/ver/example'")

func main() {
	flag.Parse()

	if *packageName == "" {
		fmt.Println("Please provide a package name with the -p parameter.")
		os.Exit(-1)
	}

	subPackages, err := getSubpackages(*packageName)

	if err != nil {
		fmt.Printf("Failed to extract package information from %s (%s) with error: %s\n", *packageName,
			os.Getenv("GOPATH")+*packageName,
			err.Error())
		os.Exit(-1)
	}

	information, err := getInformationFromPackages(subPackages)

	if err != nil {
		fmt.Printf("Failed to parse packages with error: %s\n", err.Error())
		os.Exit(-1)
	}

	for _, inf := range information {
		fmt.Println(inf)
	}
}

func getSubpackages(packageName string) ([]string, error) {
	sourceDirectory := path.Join(os.Getenv("GOPATH"), "src")

	packages := []string{}
	err := filepath.Walk(path.Join(sourceDirectory, packageName), func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			packageName := path[len(sourceDirectory)+1:]
			packages = append(packages, packageName)
		}
		return nil
	})

	return packages, err
}

func getInformationFromPackages(packageNames []string) ([]string, error) {
	information := []string{}

	for _, packageName := range packageNames {
		packageInformation, err := getInformationFromPackage(packageName)

		if err != nil {
			return information, err
		}

		for _, inf := range packageInformation {
			information = append(information, inf)
		}
	}

	return information, nil
}

func getInformationFromPackage(p string) ([]string, error) {
	var conf loader.Config
	conf.Import(p)

	prog, err := conf.Load()
	if err != nil {
		return []string{}, err
	}

	return getInformationFromProgram(p, prog), nil
}

func getInformationFromProgram(basePackage string, prog *loader.Program) []string {
	rv := []string{}

	for pkg := range prog.AllPackages {
		path := pkg.Path()

		if strings.Index(path, basePackage) == 0 {
			scope := pkg.Scope()

			for _, cs := range recurseScope(path, scope) {
				rv = append(rv, cs)
			}
		}
	}

	return rv
}

func recurseScope(path string, s *types.Scope) []string {
	rv := []string{}

	for _, sn := range s.Names() {
		lookup := s.Lookup(sn)
		lookupType := lookup.Type()

		if !lookup.Exported() {
			continue
		}

		//  Extract public fields from the structs.
		s, isStruct := lookupType.Underlying().(*types.Struct)

		if isStruct {
			msg := bytes.NewBufferString(fmt.Sprintf("%s type %s struct { ", path, sn))

			for fi := 0; fi < s.NumFields(); fi++ {
				field := s.Field(fi)

				if field.Exported() {
					//TODO: Handle nested types instead of using field.Type().String()?
					msg.WriteString(fmt.Sprintf("%s %s", field.Name(), field.Type().String()))

					if fi < s.NumFields()-2 {
						msg.WriteString(", ")
					} else {
						msg.WriteString(" ")
					}
				}
			}

			msg.WriteRune('}')

			rv = append(rv, msg.String())
			continue
		}

		rv = append(rv, fmt.Sprintf("%s %s %s", path, sn, lookupType))

		mset := types.NewMethodSet(lookupType)
		for i := 0; i < mset.Len(); i++ {
			rv = append(rv, fmt.Sprintf("%s %s %s", path, sn, mset.At(i)))
		}
	}

	return rv
}
