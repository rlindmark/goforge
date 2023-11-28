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
		{10, 0, "/first", nil, "/current", nil, 1, ""},
	}

	for _, test := range testCases {

		pagination, _ := NewPagination(test.limit, test.offset, test.first, test.previous, test.current, test.next, test.total)

		fmt.Printf("%v", pagination.asJson())
	}
	// limit := 0
	// offset := 0
	// first := "first"
	// previous := "previous"
	// current := "current"
	// next := "next"
	// total := 1

	// pagination, _ := NewPagination(limit, offset, first, &previous, current, &next, total)

	// fmt.Printf("%v", pagination.asJson())
}
