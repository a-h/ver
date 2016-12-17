package main

// PackageSignature is a map of packages to Signatures.
type PackageSignatures map[string]Signature

func NewSignature() Signature {
	return Signature{}
}

// Signature is the exported aspects of a type.
type Signature struct {
	Functions  []string
	Fields     []string
	Constants  []string
	Structs    []string
	Interfaces []string
}
