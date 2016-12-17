package example

var examplePrivateField []string

// ExamplePublicField is a public field.
var ExamplePublicField []string

type privateStruct struct {
	privateStructField string
}

type privateInterface interface {
	Close()
}

type PublicInterface interface {
	Close()
}

// PublicStructA is used to test that only public structs
// are returned.
type PublicStructA struct {
	PublicStructFieldA  string
	privateStructFieldB int
}

// PublicStructB is used to test that only public structs
// are returned.
type PublicStructB struct {
	PublicStructFieldA  string
	PublicStructFieldB  string
	privateStructFieldB int
}

func privateFunction() privateStruct {
	return privateStruct{}
}

// PublicFunctionA is used to test that only public functions
// are exported.
func PublicFunctionA(p1 string) PublicStructA {
	return PublicStructA{}
}

// PublicFunctionB is used to test that only public functions
// are exported.
func PublicFunctionB(p1 string) *PublicStructB {
	return &PublicStructB{}
}

type PublicStructC struct {
}

func (p PublicStructC) Receiver() string {
	return ""
}

func (p *PublicStructC) ReceiverPointer() string {
	return ""
}

type PublicStructD struct {
	PublicStructA
	PublicStructB
}
