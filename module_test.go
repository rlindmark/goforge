package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIsValidModuleSlug(t *testing.T) {

	var testCases = []struct {
		slug   string
		expect bool
	}{
		{"", false},
		{" ", false},
		{"0", false},
		{"A", false},
		{"a", false},
		{"a ", false},
		{" a", false},
		{"a0", false},
		{"puppetlabs-concat", true},
	}

	for _, test := range testCases {

		ok, _ := isValidModuleSlug(test.slug)

		if ok != test.expect {
			t.Errorf("For %v: got =%v, expected %v", test.slug, ok, test.expect)
		}
	}
}
func TestIsValidModuleName(t *testing.T) {

	var testCases = []struct {
		modulename string
		expect     bool
	}{
		{"", false},
		{" ", false},
		{"0", false},
		{"A", false},
		{"a", true},
		{"a ", false},
		{" a", false},
		{"a0", true},
		{"puppetlabs-concat", false},
	}

	for _, test := range testCases {

		ok, _ := isValidModuleName(test.modulename)

		if ok != test.expect {
			t.Errorf("For %v: got =%v, expected %v", test.modulename, ok, test.expect)
		}
	}
}

func TestNewModule(t *testing.T) {

	var testCases = []struct {
		modulename string
		expect     string
	}{
		{"puppetlabs-concat", ""},
		{"puppetlabs", ""},
	}

	for _, test := range testCases {

		module, err := NewModule(test.modulename)

		// module && err == nil -> pass
		// err != nil &&
		if err != nil && module != nil {
			fmt.Printf("For %v: err=%v, expected %v", test.modulename, err, test.expect)
		}
		// if err != nil
		// 	fmt.Printf("For %v: module:%v, expected %v", test.modulename, module, test.expect)
		// }
	}
}

func TestModuleMarshal(t *testing.T) {

	var testCases = []struct {
		modulename string
		expect_url string
	}{
		{"puppetlabs-module",
			`{"uri":"/v3/modules/puppetlabs-module","slug":"puppetlabs-module","name":"module","deprecated_at":null,"Owner":{"uri":"/v3/users/puppetlabs","slug":"puppetlabs","username":"puppetlabs","gravatar_id":"nogravatar"}}`,
		},
	}

	for _, test := range testCases {

		module, err := NewModule(test.modulename)
		if err != nil {
			t.Errorf("Unhandeled test: err=%v", err)
			continue
		}
		result, err := json.Marshal(module)
		if err != nil {
			t.Errorf("err:%v", err)
		}
		if string(result) != test.expect_url {
			t.Errorf("got %v, expected %v", string(result), test.expect_url)
		}
	}
}
