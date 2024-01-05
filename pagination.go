package main

import (
	"fmt"
)

type Pagination struct {
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	First    string  `json:"first"`
	Previous *string `json:"previous"`
	Current  string  `json:"current"`
	Next     *string `json:"next"`
	Total    int     `json:"total"`
}

func NewPagination(limit int, offset int, first string, previous *string, current string, next *string, total int) (Pagination, error) {
	return Pagination{limit, offset, first, previous, current, next, total}, nil
}

func (p *Pagination) asJson() string {
	result := "{"
	result += fmt.Sprintf("\"limit\":%d,", p.Limit)
	result += fmt.Sprintf("\"offset\":%d,", p.Offset)
	result += fmt.Sprintf("\"first\":%q,", p.First)
	if p.Previous == nil {
		result += "\"previous\":null,"
	} else {
		result += fmt.Sprintf("\"previous\":%v,", p.Previous)
	}
	result += fmt.Sprintf("\"current\":%q,", p.Current)
	if p.Next == nil {
		result += "\"next\":null,"
	} else {
		result += fmt.Sprintf("\"next\":%v,", p.Next)
	}
	result += fmt.Sprintf("\"total\":%d", p.Total)
	result += "}"

	return result
}
