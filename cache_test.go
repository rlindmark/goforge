package main

import (
	"fmt"
	"testing"
)

func TestFileInCache(t *testing.T) {
	var testCases = []struct {
		module_name string
		expect      bool
	}{
		{"", false}, // no empty argument
		{"p/puppetlabs/puppetlabs-stdlib-9.4.1.tar.gz", true},
		{"p/puppetlabs/puppetlabs-stdlib-9.4.0.tar.gz", true},
		{"p/puppetlabs/puppetlabs-stdlib-9.3.0.tar.gz", false},
		{"p", true},
	}

	fmt.Print("NOTE: Make sure \"make cache\" is executed before running this test")

	for _, test := range testCases {

		result := FileInCache(test.module_name)

		if result != test.expect {
			t.Errorf("TestFileInCache %v: got %v, expected %v", test.module_name, result, test.expect)
		}
	}
}

func TestModulePathInCache(t *testing.T) {
	var testCases = []struct {
		module_name string
		expect      string
	}{
		{"puppetlabs-stdlib-9.4.1.tar.gz", "cache/p/puppetlabs/puppetlabs-stdlib-9.4.1.tar.gz"},
		{"puppetlabs-stdlib-9.4.0.tar.gz", "cache/p/puppetlabs/puppetlabs-stdlib-9.4.0.tar.gz"},
		{"puppetlabs-stdlib-9.3.0.tar.gz", "cache/p/puppetlabs/puppetlabs-stdlib-9.3.0.tar.gz"},
		{"p", "cache/p/p/p"},
	}

	fmt.Print("NOTE: Make sure \"make cache\" is executed before running this test")

	for _, test := range testCases {

		result, _ := ModulePathInCache(test.module_name)

		if result != test.expect {
			t.Errorf("TestModulePathInCache %v: got %v, expected %v", test.module_name, result, test.expect)
		}
	}
}
