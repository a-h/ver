package example2

// TestField is here to test subpackages.
var TestField string

// TopLevel is used to test anonymous struct handling.
type TopLevel struct {
	Anonymous struct {
		AnonymousField  string
		DoubleAnonymous struct {
			testShouldBeMissing string
		}
	}
}
