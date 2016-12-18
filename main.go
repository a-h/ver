package main

import (
	"flag"
	"fmt"
	"go/types"
	"path/filepath"

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

func getInformationFromPackages(packageNames []string) (PackageSignatures, error) {
	var conf loader.Config

	for _, pkg := range packageNames {
		conf.Import(pkg)
	}

	prog, err := conf.Load()
	if err != nil {
		return PackageSignatures{}, err
	}

	return getInformationFromProgram(prog), nil
}

func getInformationFromProgram(prog *loader.Program) PackageSignatures {
	rv := PackageSignatures{}

	for pkg := range prog.AllPackages {
		path := pkg.Path()

		rv[path] = getSignatureFromScope(pkg.Scope())
	}

	return rv
}

func getSignatureFromScope(s *types.Scope) Signature {
	rv := NewSignature()

	for _, sn := range s.Names() {
		lookup := s.Lookup(sn)
		lookupType := lookup.Type()

		if !lookup.Exported() {
			continue
		}

		switch lookup.(type) {
		case *types.Func:
			rv.Functions = append(rv.Functions, lookup.String())
			break
		case *types.Var:
			rv.Fields = append(rv.Fields, lookup.String())
			break
		case *types.Const:
			value := lookup.(*types.Const).Val().String()
			rv.Constants = append(rv.Constants, lookup.String()+" = "+value)
			break
		}

		switch lookupType.Underlying().(type) {
		case *types.Struct:
			rv.Structs = append(rv.Structs, renderStruct(sn, lookupType.Underlying().(*types.Struct)))
			break
		case *types.Interface:
			rv.Interfaces = append(rv.Interfaces, lookupType.String())
			break
		}

		// Extract methods from structs, interfaces and pointers to structs.
		for _, msetType := range []types.Type{lookupType, types.NewPointer(lookupType)} {
			mset := types.NewMethodSet(msetType)
			for i := 0; i < mset.Len(); i++ {
				method := mset.At(i)
				if method.Obj().Exported() {
					rv.Functions = append(rv.Functions, method.String())
				}
			}
		}
	}

	return rv
}

func renderStruct(name string, s *types.Struct) string {
	msg := bytes.NewBufferString("struct")

	if name != "" {
		msg.WriteString(" " + name)
	}

	msg.WriteString(" {")

	fieldCount := s.NumFields()

	for fi := 0; fi < fieldCount; fi++ {
		field := s.Field(fi)

		if field.Exported() {
			s, isStruct := field.Type().Underlying().(*types.Struct)

			if isStruct {
				msg.WriteString(fmt.Sprintf(" %s %s", field.Name(), renderStruct("", s)))
			} else {
				msg.WriteString(" " + field.String())
			}

			if fi < fieldCount-2 {
				msg.WriteString(",")
			} else {
				msg.WriteString(" ")
			}
		}
	}

	msg.WriteString("}")

	return msg.String()
}
