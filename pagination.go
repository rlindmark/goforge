package main

import (
	"fmt"
	"net/url"
	"strconv"
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

func CreatePagination(query url.Values, total int) (*Pagination, error) {
	var previous_page *string
	var current_page *string
	var next_page *string

	m := query

	url_limit := query.Get("limit")
	if url_limit == "" {
		return nil, fmt.Errorf("limit not defined in query %+v", query)
	}

	limit, err := strconv.Atoi(url_limit)
	if err != nil {
		return nil, fmt.Errorf("limit=%v is not an interger", url_limit)
	}

	url_offset := query.Get("offset")
	if url_offset == "" {
		return nil, fmt.Errorf("offset not defined in query %+v", query)
	}

	offset, err := strconv.Atoi(url_offset)
	if err != nil {
		return nil, fmt.Errorf("offset=%v is not an interger", url_offset)
	}

	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("limit=%d must be between 1..100", limit)
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset=%d must be >= 0", offset)
	}
	if total < 0 {
		return nil, fmt.Errorf("total=%d must be >= 0", total)
	}

	first_page := fmt.Sprintf("/v3/releases?offset=%d&limit=%d", offset, limit)

	//"first": "/v3/modules?offset=0&limit=20",
	if offset < limit {
		previous_page = nil
	} else {
		// previous_offset = offset - limit, limit unchanged
		previous_query := m
		previous_query.Set("offset", fmt.Sprint(offset-limit))
		previous_escaped, _ := url.QueryUnescape(previous_query.Encode())
		previous_unescaped, _ := url.QueryUnescape(previous_escaped)
		str := "/v3/releases?" + previous_unescaped
		previous_page = &str
	}

	current_query := m
	current_query.Set("offset", fmt.Sprint(offset))
	current_escaped, _ := url.QueryUnescape(current_query.Encode())
	current_unescaped, _ := url.QueryUnescape(current_escaped)
	str := "/v3/releases?" + current_unescaped
	current_page = &str

	if offset+limit < total {
		next_query := m
		next_query.Set("offset", fmt.Sprint(offset+limit))
		next_escaped, _ := url.QueryUnescape(next_query.Encode())
		next_unescaped, _ := url.QueryUnescape(next_escaped)
		str2 := "/v3/releases?" + next_unescaped
		next_page = &str2
	} else {
		next_page = nil
	}

	pagination := Pagination{
		Limit:    limit,
		Offset:   offset,
		First:    first_page,
		Previous: previous_page,
		Current:  *current_page,
		Next:     next_page,
		Total:    total,
	}

	return &pagination, nil
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
