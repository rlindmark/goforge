package main

import (
	"fmt"
	"regexp"
)

type Owner struct {
	Uri         string `json:"uri"` // /v3/users/puppetlabs
	Slug        string `json:"slug"`
	Username    string `json:"username"`
	Gravatar_id string `json:"gravatar_id"`
}

func isValidOwnerSlug(slug string) (bool, error) {
	return regexp.MatchString("^[a-zA-Z0-9]+$", slug)
}

func NewOwner(uri string, slug string, username string, gravatar_id string) (*Owner, error) {
	ok, _ := isValidOwnerSlug(slug)
	if !ok {
		return nil, fmt.Errorf("cant create owner with invalid slug: %v", slug)
	}
	return &Owner{uri, slug, username, gravatar_id}, nil
}
