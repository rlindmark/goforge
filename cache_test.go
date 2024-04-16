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
		{"cache/p/puppetlabs/puppetlabs-stdlib-9.3.1.tar.gz", false},
		{"p", false},
	}

	for _, test := range testCases {

		result := FileInCache(test.module_name)

		if result != test.expect {
			t.Errorf("TestFileInCache module:%v, got:%v, expected:%v", test.module_name, result, test.expect)
		}
	}
}

// func TestModulePathInCache(t *testing.T) {
// 	var testCases = []struct {
// 		module_name string
// 		expect      string
// 	}{
// 		{"puppetlabs-stdlib-9.4.1.tar.gz", "cache/p/puppetlabs/puppetlabs-stdlib-9.4.1.tar.gz"},
// 		{"puppetlabs-stdlib-9.4.0.tar.gz", "cache/p/puppetlabs/puppetlabs-stdlib-9.4.0.tar.gz"},
// 		{"puppetlabs-stdlib-9.3.0.tar.gz", "cache/p/puppetlabs/puppetlabs-stdlib-9.3.0.tar.gz"},
// 		{"p", "cache/p/p/p"},
// 	}

// 	for _, test := range testCases {

// 		result, _ := ModulePathInCache(test.module_name)

// 		if result != test.expect {
// 			t.Errorf("TestModulePathInCache %v: got %v, expected %v", test.module_name, result, test.expect)
// 		}
// 	}
// }

func TestGetModuleFilePath(t *testing.T) {
	var testCases = []struct {
		module_name string
		expect      string
	}{
		{"", ""},
		{"p", ""},
		{"puppetlabs-stdlib-9.4.1.tar.gz", ""},
		{"puppetlabs-stdlib-9.4.0", "cache/p/puppetlabs/puppetlabs-stdlib-9.4.0.tar.gz"},
		{"puppetlabs-stdlib-9.3.0", "cache/p/puppetlabs/puppetlabs-stdlib-9.3.0.tar.gz"},
	}

	c := NewForgeCache("cache")

	for _, test := range testCases {

		result, err := c.GetModuleFilePath(test.module_name)

		if result != test.expect {
			t.Errorf("GetModuleFilePath(%v): got %v err=%v, expected %v", test.module_name, result, err, test.expect)
		}
	}
}

func TestGetModuleVersions(t *testing.T) {
	var testCases = []struct {
		module_name string
		expect      []string
	}{
		{"", []string{""}},
		{"p", []string{""}},
		{"puppetlabs-stdlib", []string{""}},
		{"puppetlabs-stdlib-9.4.0", []string{"cache/p/puppetlabs/puppetlabs-stdlib-9.4.0.tar.gz"}},
		// {"puppetlabs-stdlib-9.3.0", []string{"cache/p/puppetlabs/puppetlabs-stdlib-9.3.0.tar.gz"}},
	}

	c := NewForgeCache("cache")

	for _, test := range testCases {

		result := c.GetModuleVersions(test.module_name)

		fmt.Printf("GetModuleVersions(%v)=%v\n", test.module_name, result)
		// if result != test.expect {
		// 	t.Errorf("GetModuleFilePath(%v): got %v err=%v, expected %v", test.module_name, result, err, test.expect)
		// }
	}
}

func TestGetAllUsers(t *testing.T) {

	c := NewForgeCache("cache")

	result := c.GetAllUsers()

	if len(result) != 2 {
		t.Errorf("expected [\"pdxcat\",\"puppetlabs\"] got %v", result)
	}
}
