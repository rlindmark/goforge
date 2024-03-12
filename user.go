package main

import (
	"fmt"
	"regexp"
)

type User struct {
	Uri           string `json:"uri"`           // <URL path and query string> Relative URL for this User resource
	Slug          string `json:"slug"`          // ^[a-zA-Z0-9]+$ Unique textual identifier for this User resource
	Gravatar_Id   string `json:"gravatar_id"`   // Gravatar ID, learn more at https://gravatar.com/
	Username      string `json:"username"`      // Username for this User resource
	Display_name  string `json:"display_name"`  // Free form display name for this User resource
	Release_Count int    `json:"release_count"` // Total number of module releases (versions) published under this User
	Module_Count  int    `json:"module_count"`  // Total number of unique modules published under this User
	Created_At    string `json:"created_at"`    // <iso8601> Date and time this User resource was created
	Updated_At    string `json:"updated_at"`    // <iso8601> Date and time this User resource was last modified
}

func IsValidUserSlug(user_slug string) (bool, error) {
	return regexp.MatchString("^[a-zA-Z0-9]+$", user_slug)
}

func NewUser(
	uri string, slug string, gravatar_id string, username string,
	display_name string, release_count int, module_count int,
	created_at string, updated_at string) (*User, error) {

	ok, _ := IsValidUserSlug(slug)
	if !ok {
		return nil, fmt.Errorf(`{"message":"400 Bad Request","errors":["'%s' is not a valid user slug"]}`, slug)
	}

	// FIXME: assert that uri begins with "/v3/users/"
	// FIXME: gravatar_id should be sha256 char long.
	// FIXME: Need to check that the user exist in cache
	// FIXME: get release_count
	// FIXME: get module_count
	return &User{uri, slug, gravatar_id, username, display_name, release_count, module_count, created_at, updated_at}, nil
}
