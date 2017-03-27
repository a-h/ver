package main

import "fmt"

// Version represents a major, minor and build version.
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Build int `json:"build"`
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Build)
}

// MarshalJSON determines the output format of the version struct in JSON.
func (v Version) MarshalJSON() ([]byte, error) {
	return []byte("\"" + v.String() + "\""), nil
}

// Add adds two versions together.
func (v Version) Add(d Version) Version {
	return Version{
		Major: v.Major + d.Major,
		Minor: v.Minor + d.Minor,
		Build: v.Build + d.Build,
	}
}
