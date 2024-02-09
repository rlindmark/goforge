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
		{"cache/p/puppetlabs/puppetlabs-stdlib-9.4.1.tar.gz", true},
		{"cache/p/puppetlabs/puppetlabs-stdlib-9.4.0.tar.gz", true},
		{"cache/p/puppetlabs/puppetlabs-stdlib-9.3.0.tar.gz", false},
		{"p", false},
	}

	for _, test := range testCases {

		result := FileInCache(test.module_name)

		if result != test.expect {
			t.Errorf("TestFileInCache module:%v, got:%v, expected:%v", test.module_name, result, test.expect)
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

	for _, test := range testCases {

		result, _ := ModulePathInCache(test.module_name)

		if result != test.expect {
			t.Errorf("TestModulePathInCache %v: got %v, expected %v", test.module_name, result, test.expect)
		}
	}
}
