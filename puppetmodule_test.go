package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewPuppetModule(t *testing.T) {

	var testCases = []struct {
		uri    string
		expect bool
	}{
		{"puppetlabs-stdlib-1.0.0", false},
	}

	for _, test := range testCases {

		result, err := NewPuppetModule(test.uri)

		if result == nil && test.expect != false {
			t.Errorf("For %v: got %v err=%v, expected %v", test.uri, result, err, test.expect)
		}
		if result != nil && test.expect != true {
			t.Errorf("For %v: got %v err=%v, expected %v", test.uri, result, err, test.expect)
		}
	}
}

func errorf400(module string) error {
	return fmt.Errorf(`{"message":"400 Bad Request","errors":["'%s' is not a valid release slug"]}`, module)
}

func TestValidModuleReleaseFile(t *testing.T) {
	var testCases = []struct {
		uri    string
		err    error
		expect bool
	}{
		{"puppetlabs-stdlib-1.0.0.tar.gz", nil, true},
		{"puppetlabs-stdlib-1.0.0.tar", errorf400("puppetlabs-stdlib-1.0.0.tar"), false},
		{"puppetlabs-stdlib-1.0.0.gz", errorf400("puppetlabs-stdlib-1.0.0.gz"), false},
		{"", errorf400(""), false},
		{"tar.gz", errorf400("tar.gz"), false},
		{"../puppetlabs-stdlib-1.0.0.tar.gz", errorf400("../puppetlabs-stdlib-1.0.0.tar.gz"), false},
		{"/puppetlabs-stdlib-1.0.0.tar.gz", errorf400("/puppetlabs-stdlib-1.0.0.tar.gz"), false},
	}

	for _, test := range testCases {

		result, err := validModuleReleaseFilename(test.uri)
		if result == test.expect {
			if err == nil && test.err == nil {
				fmt.Printf("pass\n")
			} else {
				//if !errors.Is(err, test.err) {
				if err.Error() != test.err.Error() {
					t.Errorf("expected %t, got %t, err = %v, wants err = %v\n", test.expect, result, err, test.err)
				}
			}
		} else {
			if err != test.err {
				t.Errorf("expected %t, got %t, err = %v, wants err = %v\n", test.expect, result, err, test.err)
			}
		}
		// 	fmt.Printf("result = nil and err=%s\n", err)
		// }
	}

}
func TestPuppetModuleMarshal(t *testing.T) {

	var testCases = []struct {
		uri    string
		expect bool
	}{
		{"puppetlabs-stdlib-9.4.0", true},
		{"puppetlabs-stdlib-1.0.0", false},
	}

	for _, test := range testCases {

		result, err := NewPuppetModule(test.uri)
		if err != nil && test.expect == true {
			t.Errorf("cant create module %s, err = %s", test.uri, err)
		}
		jSON, _ := json.Marshal(result)
		if jSON != nil {
			fmt.Print(string(jSON))
		}
	}
}

func TestGet_metadata(t *testing.T) {
	var testCases = []struct {
		module string
		expect bool
	}{
		{"pdxcat-nrpe-2.1.1", true},
		{"puppetlabs-stdlib-9.4.1", true},
	}

	for _, test := range testCases {

		json, err := get_metadata(test.module)
		if err != nil {
			t.Errorf("cant get metadata for module %s, err = %s", test.module, err)
		}
		if json != "" {
			fmt.Printf("json:%v", string(json))
		}
	}
}
