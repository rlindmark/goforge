package main

import (
	"fmt"
	"testing"
)

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
		{10, 0, "/first", nil, "/current", nil, 1,
			"{\"limit\":10,\"offset\":0,\"first\":\"/first\",\"previous\":null,\"current\":\"/current\",\"next\":null,\"total\":1}"},
	}

	for _, test := range testCases {

		pagination, _ := NewPagination(test.limit, test.offset, test.first, test.previous, test.current, test.next, test.total)

		result := pagination.asJson()
		if result != test.expect {
			fmt.Printf("Got %v, expected %v", result, test.expect)
		}
	}
}
