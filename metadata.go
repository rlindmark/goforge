package main

type Metadata struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Author       string `json:"author"`
	License      string `json:"license"`
	Summary      string `json:"summary"`
	Source       string `json:"source"`
	Project_page string `json:"project_page"`

	Issues_url   string       `json:"issues_url"`
	Dependencies []Dependency `json:"dependencies"`

	Data_provider string `json:"data_provider"`

	Operatingsystem_support []OperatingSystemSupport `json:"operatingsystem_support"`
	Requirements            []Requirement            `json:"requirements"`
}

type OperatingSystemSupport struct {
	Operatingsystem        string   `json:"operatingsystem"`
	Operatingsystemrelease []string `json:"operatingsystemrelease"`
}

type Requirement struct {
	Name                string `json:"name"`
	Version_requirement string `json:"version_requirement"`
}

type Dependency struct {
	Name                string `json:"name"`
	Version_requirement string `json:"version_requirement"`
}
