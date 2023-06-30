package main

import (
	"fmt"
	"testing"
)

func TestLoadFromJSON(t *testing.T) {
	data, err := loadFromJSON("/Users/lodrem/w/typedb/london-tube-data/loader/train-network.json")
	if err != nil {
		t.Fatalf("no error: %s", err)
	}

	fmt.Printf("stations: %v\n", data.Stations)
	fmt.Printf("lines: %v\n", data.Lines)
}
