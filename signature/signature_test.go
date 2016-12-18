package signature

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
		expected Signature
	}{
		{
			name: "Public functions are extracted",
			code: []string{"package nonexistent", "func GetStringPackageFunction() string { return \"Hello\" }"},
			expected: Signature{
				Functions: []string{"func github.com/a-h/nonexistent.GetStringPackageFunction() string"},
			},
		},
		{
			name:     "Private functions are not extracted",
			code:     []string{"package nonexistent", "func getStringPackageFunction() string { return \"Hello\" }"},
			expected: Signature{},
		},
		{
			name: "Public structs are extracted, and only their public fields are exported",
			code: []string{"package nonexistent", "type Test struct { Public string; private string; Public2 string }"},
			expected: Signature{
				Structs: []string{"struct Test { field Public string, field Public2 string }"},
			},
		},
		{
			name: "Public receiver methods are extracted",
			code: []string{"package nonexistent", "type Test struct { value string }", "func (t Test) GetStringReceiver() string { return t.value }"},
			expected: Signature{
				Structs: []string{"struct Test {}"},
				Functions: []string{
					"method (github.com/a-h/nonexistent.Test) GetStringReceiver() string",
					"method (*github.com/a-h/nonexistent.Test) GetStringReceiver() string",
				},
			},
		},
		{
			name: "Public pointer receiver methods are extracted",
			code: []string{"package nonexistent", "type Test struct { value string }", "func (t *Test) GetStringPointerReceiver() string { return t.value }"},
			expected: Signature{
				Structs:   []string{"struct Test {}"},
				Functions: []string{"method (*github.com/a-h/nonexistent.Test) GetStringPointerReceiver() string"},
			},
		},
		{
			name: "Private receiver methods are not extracted",
			code: []string{"package nonexistent", "type Test struct { value string }", "func (t Test) getString() string { return t.value }"},
			expected: Signature{
				Structs: []string{"struct Test {}"},
			},
		},
		{
			name: "Package level fields should be extracted",
			code: []string{"package nonexistent", "var Public int"},
			expected: Signature{
				Fields: []string{"var github.com/a-h/nonexistent.Public int"},
			},
		},
		{
			name:     "Private package level fields are not extracted",
			code:     []string{"package nonexistent", "var private int"},
			expected: Signature{},
		},
		{
			name: "Public constants should be extracted and should include their value",
			code: []string{"package nonexistent", "const HTTPNotFound = 400"},
			expected: Signature{
				Constants: []string{"const github.com/a-h/nonexistent.HTTPNotFound untyped int = 400"},
			},
		},
		{
			name:     "Private constants should not be extracted",
			code:     []string{"package nonexistent", "const httpNotFound = 400"},
			expected: Signature{},
		},
		{
			name: "Public interfaces should be extracted",
			code: []string{"package nonexistent", "type Test interface { Close() }"},
			expected: Signature{
				Functions:  []string{"method (github.com/a-h/nonexistent.Test) Close()"},
				Interfaces: []string{"github.com/a-h/nonexistent.Test"},
			},
		},
		{
			name:     "Types are extracted",
			code:     []string{"package nonexistent", "type Test int"},
			expected: Signature{},
		},
		{
			name: "Anonymous nested structs are extracted without public fields",
			code: []string{"package nonexistent", "type Test struct { A struct{B string; c string} }"},
			expected: Signature{
				Structs: []string{"struct Test { A struct { field B string } }"},
			},
		},
	}
	for _, tt := range tests {
		pkg, err := parseGoIntoPackage(basePackage, strings.Join(tt.code, "\n"))

		if err != nil {
			t.Errorf("%s - failed to parse Go with error %v", tt.name, err)
			continue
		}

		actual := GetFromScope(pkg.Scope())

		compareLengths(tt.name, "Functions", tt.expected.Functions, actual.Functions, t)
		compareElements(tt.name, "Functions", tt.expected.Functions, actual.Functions, t)

		compareLengths(tt.name, "Fields", tt.expected.Fields, actual.Fields, t)
		compareElements(tt.name, "Fields", tt.expected.Fields, actual.Fields, t)

		compareLengths(tt.name, "Constants", tt.expected.Constants, actual.Constants, t)
		compareElements(tt.name, "Constants", tt.expected.Constants, actual.Constants, t)

		compareLengths(tt.name, "Structs", tt.expected.Structs, actual.Structs, t)
		compareElements(tt.name, "Structs", tt.expected.Structs, actual.Structs, t)

		compareLengths(tt.name, "Interfaces", tt.expected.Interfaces, actual.Interfaces, t)
		compareElements(tt.name, "Interfaces", tt.expected.Interfaces, actual.Interfaces, t)
	}
}

func compareElements(testname string, element string, expected []string, actual []string, t *testing.T) {
	max := len(actual)
	if max < len(expected) {
		max = len(expected)
	}

	var buf bytes.Buffer

	errorOccurred := false
	for i := 0; i < max; i++ {
		expectedItem := getItemOrDefault(expected, i, "<missing>")
		actualItem := getItemOrDefault(actual, i, "<missing>")

		buf.WriteString(fmt.Sprintf("%s - %d - expected '%s', but got '%s'\n",
			testname, i, expectedItem, actualItem))

		if actualItem != expectedItem {
			errorOccurred = true
		}
	}

	if errorOccurred {
		t.Error(buf.String())
	}
}

func compareLengths(testname string, element string, expected []string, actual []string, t *testing.T) {
	if len(actual) != len(expected) {
		t.Errorf("%s - expected %d extracted, but %d were extracted", testname, len(expected), len(actual))
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
