package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"
)

func Test_GetInformationFromProgram(t *testing.T) {
	basePackage := "github.com/a-h/nonexistent"

	tests := []struct {
		name     string
		code     []string
		expected []string
	}{
		{
			name:     "Package level struct",
			code:     []string{"package nonexistent", "type X struct { Field string }"},
			expected: []string{"github.com/a-h/nonexistent type X struct { Field string }"},
		},
		{
			name:     "Package level field",
			code:     []string{"package nonexistent", "type S string"},
			expected: []string{"github.com/a-h/nonexistent S github.com/a-h/nonexistent.S"},
		},
	}
	for _, tt := range tests {
		pkg, err := parseGoIntoPackage(basePackage, strings.Join(tt.code, "\n"))

		if err != nil {
			t.Errorf("%s - failed to parse Go with error %v", tt.name, err)
			continue
		}

		actual := recurseScope(basePackage, pkg.Scope())

		if len(actual) != len(tt.expected) {
			t.Errorf("%s - expected %d extracted, but %d were extracted", tt.name, len(tt.expected), len(actual))
		}

		max := len(actual)
		if max < len(tt.expected) {
			max = len(tt.expected)
		}

		var buf bytes.Buffer

		errorOccurred := false
		for i := 0; i < max; i++ {
			actualLine := getItemOrDefault(actual, i, "<missing>")
			expectedLine := getItemOrDefault(tt.expected, i, "<missing>")

			buf.WriteString(fmt.Sprintf("%s - %d - expected '%s', but got '%s'\n",
				tt.name, i, expectedLine, actualLine))

			if actualLine != expectedLine {
				errorOccurred = true
			}
		}

		if errorOccurred {
			t.Error(buf.String())
		}
	}
}

func getItemOrDefault(items []string, index int, missing string) string {
	if index >= len(items) {
		return missing
	}

	return items[index]
}

func parseGoIntoPackage(packageName string, input string) (*types.Package, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", input, 0)
	if err != nil {
		return nil, err
	}

	var conf types.Config
	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	return conf.Check(packageName, fset, []*ast.File{f}, &info)
}
