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
		{"puppet labs", false},
		{"", false},
	}

	for _, test := range testCases {

		result, _ := isValidOwnerSlug(test.slug)

		if result != test.expect {
			fmt.Printf("For %v: got %v, expected %v", test.slug, result, test.expect)
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

func TestOwnerAsJson(t *testing.T) {

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
	}

	for _, test := range testCases {

		owner, _ := NewOwner(test.uri, test.slug, test.username, test.gravatar_id)

		//result := user.asJson()
		result, _ := json.Marshal(owner)
		if string(result) != test.expect {
			fmt.Printf("TestOwnerAsJson:Got %v, expected %v", result, test.expect)
		}
	}
}
