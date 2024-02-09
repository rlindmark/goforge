package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestMetadata(t *testing.T) {
	filename := "tests/metadata.json"
	jsonFile, err := os.Open(filename)
	if err != nil {
		t.Errorf("Unable to open test file %s", filename)
	}

	var metadata Metadata
	bytesValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(bytesValue, &metadata)
	//fmt.Println(metadata.asJSON())
	json, _ := json.Marshal(metadata)
	fmt.Println(string(json))
}
