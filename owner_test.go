package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIsValidOwnerSlug(t *testing.T) {

	var testCases = []struct {
		slug   string
		expect bool
	}{
		{"puppetlabs", true},
		{"Puppetlabs", true},
		{"p", true},
		{"0", true},
		{"0a", true},
		{"puppet labs", false},
		{"", false},
		{"ö", false},
		{"aö", false},
	}

	for _, test := range testCases {

		result, _ := isValidOwnerSlug(test.slug)

		if result != test.expect {
			t.Errorf("Testing %v: got %v, expected %v", test.slug, result, test.expect)
		}
	}
}

func TestNewOwner(t *testing.T) {

	var testCases = []struct {
		uri         string
		slug        string
		username    string
		gravatar_id string
		expect      bool
	}{
		{"puppetlabs", "puppetlabs", "puppetlabs", "gravatar", true},
		{"puppetlabs", "", "puppetlabs", "gravatar", false},
		{"puppetlabs", "puppet labs", "puppetlabs", "gravatar", false},
	}

	for _, test := range testCases {

		result, err := NewOwner(test.uri, test.slug, test.username, test.gravatar_id)

		if result == nil && test.expect != false {
			fmt.Printf("For %v: got %v err=%v, expected %v", test.slug, result, err, test.expect)
		}
		if result != nil && test.expect != true {
			fmt.Printf("For %v: got %v err=%v, expected %v", test.slug, result, err, test.expect)
		}
	}
}

func TestOwnerMarshal(t *testing.T) {

	var testCases = []struct {
		uri         string
		slug        string
		username    string
		gravatar_id string
		expect      string
	}{
		{"Puppetlabs", "Slug", "Username", "Gravatar",
			`{"uri":"Puppetlabs","slug":"Slug","username":"Username","gravatar_id":"Gravatar"}`,
		},
		{"module", "slug", "username", "gravatar",
			`{"uri":"module","slug":"slug","username":"username","gravatar_id":"gravatar"}`,
		},
	}

	for _, test := range testCases {

		owner, _ := NewOwner(test.uri, test.slug, test.username, test.gravatar_id)

		result, _ := json.Marshal(owner)
		if string(result) != test.expect {
			t.Errorf("got %v, expected %v", string(result), test.expect)
		}
	}
}
