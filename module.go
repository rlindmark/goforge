package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Module struct {
	Uri           string  `json:"uri"`
	Slug          string  `json:"slug"`
	Name          string  `json:"name"`
	Deprecated_at *string `json:"deprecated_at"`
	owner         *Owner
}

// module_name is the string "owner-module" with no version
func NewModule(owner_module string) (*Module, error) {
	// uri = /v3/releases/puppetlabs-apache-4.0.0
	//ok, _ := isValidModuleName(module_name)

	module_uri := fmt.Sprintf("/v3/modules/%s", owner_module)
	module_slug := owner_module
	module_name := strings.Split(owner_module, "-")[1]
	module_owner := strings.Split(owner_module, "-")[0]
	var module_deprecated_at *string = nil

	owner_uri := fmt.Sprintf("/v3/users/%s", module_owner)
	owner_slug := module_owner
	owner_username := module_owner
	owner_gravatar_id := "nogravatar"
	owner, _ := NewOwner(owner_uri, owner_slug, owner_username, owner_gravatar_id)

	// if !ok {
	// 	return nil, fmt.Errorf("cant create module with invalid uri: %v", uri)
	// }
	module := Module{module_uri, module_slug, module_name, module_deprecated_at, owner}
	return &module, nil
}

func (module *Module) asJson() string {
	jsons := "{"
	jsons += fmt.Sprintf("%q:%q,", "uri", module.Uri)
	jsons += fmt.Sprintf("%q:%q,", "slug", module.Slug)
	jsons += fmt.Sprintf("%q:%q,", "name", module.Name)
	jsons += fmt.Sprintf("%q:null,", "deprecated_at")
	owner_jsons, _ := json.Marshal(module.owner)
	jsons += fmt.Sprintf("%q:%s", "owner", string(owner_jsons))
	//json += fmt.Sprintf("%q:%s", "owner", module.owner.asJson())
	jsons += "}"
	return jsons
}
