package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestCreatePagination(t *testing.T) {

	var testCases = []struct {
		params string // map[string]string
		total  int
		expect string
		err    error
	}{
		{ // limit not defined
			"",
			1,
			"",
			nil,
		},
		{ // limit not an integer
			"limit=a",
			1,
			"",
			nil,
		},
		{ // offset not defined
			"limit=1",
			1,
			"",
			nil,
		},
		{ // offset not an integer
			"limit=1&offset=a",
			1,
			"",
			nil,
		},
		{ // limit between 0..100
			"limit=-1&offset=0",
			1,
			"",
			nil,
		},
		{ // limit between 0..100
			"limit=101&offset=0",
			1,
			"",
			nil,
		},
		{ // total > 0
			"limit=1&offset=0",
			-1,
			"",
			nil,
		},
		{
			"offset=0&name=puppetlabs-stdlib-9.0.1",
			1,
			"",
			nil,
		},
		{
			"limit=1&name=puppetlabs-stdlib-9.0.1",
			1,
			"",
			nil,
		},
		{
			"limit=20&offset=0&name=puppetlabs-stdlib-9.0.1",
			1,
			`{"limit":20,"offset":0,"first":"/v3/releases?offset=0&limit=20","previous":null,"current":"/v3/releases?limit=20&name=puppetlabs-stdlib-9.0.1&offset=0","next":null,"total":1}`,
			nil,
		},
	}

	for _, test := range testCases {

		params, err := url.ParseQuery(test.params)
		if err != nil {
			t.Fatalf("Incorrect test data")
		}
		pagination, err := CreatePagination(params, test.total)

		if err != nil {
			fmt.Printf("error %v\n", err)
		} else {
			result := pagination.asJSON()
			if result != test.expect {
				fmt.Printf("Got %v, expected %v\n", result, test.expect)
			}
		}
	}
}

func TestNewPagination(t *testing.T) {

	var testCases = []struct {
		limit    int
		offset   int
		first    string
		previous *string
		current  string
		next     *string
		total    int
		expect   string
	}{
		{10, 0, "first", nil, "/current", nil, 1,
			`{"limit":10,"offset":0,"first":"/first","previous":null,"current":"/current","next":null,"total":1}`},
	}

	for _, test := range testCases {

		pagination, _ := NewPagination(test.limit, test.offset, test.first, test.previous, test.current, test.next, test.total)

		result := pagination.asJSON()
		if result != test.expect {
			fmt.Printf("Got %v, expected %v\n", result, test.expect)
		}
	}
}
