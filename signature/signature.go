package signature

import (
	"bytes"
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"
)

// PackageSignatures is a map of packages to Signatures.
type PackageSignatures map[string]Signature

// NewSignature creates an empty signature struct.
func NewSignature() Signature {
	return Signature{}
}

// Signature is the exported aspects of a type.
type Signature struct {
	Functions  []string `json:"functions"`
	Fields     []string `json:"fields"`
	Constants  []string `json:"constants"`
	Structs    []string `json:"structs"`
	Interfaces []string `json:"interfaces"`
}

// GetFromDirectory gets the signature of a directory of Go files.
func GetFromDirectory(dir string) (PackageSignatures, error) {
	// Iterate subdirectories too.
	directories, err := walkDirectories(dir)

	if err != nil {
		return PackageSignatures{}, err
	}

	// Import the directories
	var conf loader.Config

	for _, d := range directories {
		filenames, err := getFiles(d)

		if err != nil {
			return PackageSignatures{}, err
		}

		conf.CreatePkgs = append(conf.CreatePkgs, loader.PkgSpec{Path: d, Filenames: filenames})
	}

	prog, err := conf.Load()

	if err != nil {
		return PackageSignatures{}, err
	}

	return GetFromProgram(prog, dir), err
}

func walkDirectories(dir string) ([]string, error) {
	rv := []string{}

	err := filepath.Walk(dir, func(walkedPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() && f.Name()[0] == '.' {
			return filepath.SkipDir
		}

		if !f.IsDir() {
			return nil
		}

		rv = append(rv, walkedPath)

		return nil
	})

	return rv, err
}

func getFiles(dir string) ([]string, error) {
	files := []string{}

	fi, err := ioutil.ReadDir(dir)

	if err != nil {
		return files, err
	}

	for _, f := range fi {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".go" {
			files = append(files, path.Join(dir, f.Name()))
		}
	}

	return files, nil
}

// GetFromProgram gets a set of signatures for a program loaded with the loader.Config.
// Only packages with a matching prefix will be extracted.
func GetFromProgram(prog *loader.Program, prefix string) PackageSignatures {
	rv := PackageSignatures{}

	for pkg := range prog.AllPackages {
		path := pkg.Path()

		// Filter by prefix.
		if strings.Index(path, prefix) != 0 {
			continue
		}

		rv[path] = GetFromScope(pkg.Scope())
	}

	return rv
}

// GetFromScope gets a Signature for a given Scope.
func GetFromScope(s *types.Scope) Signature {
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
