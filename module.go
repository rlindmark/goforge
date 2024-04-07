package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Module struct {
	Uri           string  `json:"uri"`
	Slug          string  `json:"slug"`
	Name          string  `json:"name"`
	Deprecated_at *string `json:"deprecated_at"`
	Owner         *Owner  `json:"owner"`
}

func isValidModuleSlug(slug string) (bool, error) {
	// ^[a-zA-Z0-9]+[-\/][a-z][a-z0-9_]*$
	return regexp.MatchString("^[a-zA-Z0-9]+[-/][a-z][a-z0-9_]*$", slug)
}

func isValidModuleName(name string) (bool, error) {
	// ^[a-z][a-z0-9_]*$
	return regexp.MatchString("^[a-z][a-z0-9_]*$", name)
}

// module_name is the string "owner-module" with no version
func NewModule(owner_module string) (*Module, error) {
	// uri = /v3/releases/puppetlabs-apache-4.0.0
	//ok, _ := isValidModuleName(module_name)

	ok, _ := isValidModuleSlug(owner_module)
	if !ok {
		return nil, fmt.Errorf("not a valid module slug:%v", owner_module)
	}
	// isValidModuleName ^[a-z][a-z0-9_]*$
	module_uri := fmt.Sprintf("/v3/modules/%s", owner_module)
	module_slug := owner_module
	// FIXME: need to validate owner_module
	module_name := strings.Split(owner_module, "-")[1]
	module_owner := strings.Split(owner_module, "-")[0]

	var module_deprecated_at *string = nil

	owner_uri := fmt.Sprintf("/v3/users/%s", module_owner)
	owner_slug := module_owner
	owner_username := module_owner
	owner_gravatar_id := GetGravatarId(owner_slug)
	owner, err := NewOwner(owner_uri, owner_slug, owner_username, owner_gravatar_id)
	if err != nil {
		return nil, fmt.Errorf("cant create module with invalid uri:%v. Err=%v", owner_uri, err)
	}

	module := Module{module_uri, module_slug, module_name, module_deprecated_at, owner}
	return &module, nil
}
