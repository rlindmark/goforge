package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestMetadata(t *testing.T) {
	var testCases = []struct {
		metadata_file string
		expect        bool
	}{
		{"testdata/metadata.json", true},
		{"testdata/metadata1.json", true},
	}

	for _, test := range testCases {

		filename := test.metadata_file
		jsonFile, err := os.Open(filename)

		if err != nil {
			t.Errorf("Unable to open test file %s", filename)
			continue
		}

		defer jsonFile.Close()

		var metadata Metadata
		bytesValue, _ := io.ReadAll(jsonFile)
		err = json.Unmarshal(bytesValue, &metadata)
		if err != nil {
			t.Errorf("unable to unmarshal %v\n", filename)
		}
		json, _ := json.Marshal(metadata)
		fmt.Println(string(json))
		// if string(bytesValue) != string(json) {
		// 	t.Errorf("bytes:%v got:%v", string(bytesValue), string(json))
		// }
	}
}
