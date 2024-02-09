package main

import "fmt"

type Metadata struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Author       string `json:"author"`
	Summary      string `json:"summary"`
	License      string `json:"license"`
	Source       string `json:"source"`
	Project_page string `json:"project_page"`

	Issues_url   string       `json:"issues"`
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

func (m *OperatingSystemSupport) asJSON() string {
	json := "{"
	json += fmt.Sprintf("\"operatingsystem\":\"%s\"", m.Operatingsystem)

	operatingsystemreleases := m.Operatingsystemrelease
	size := len(operatingsystemreleases)
	if size > 0 {
		json += ",\"operatingsystemrelease\":["
		operatingsystemreleases := m.Operatingsystemrelease
		size := len(operatingsystemreleases)
		for index, operatingsystemrelease := range operatingsystemreleases {
			json += fmt.Sprintf("\"%v\"", operatingsystemrelease)
			if index < size-1 {
				json += ","
			}
		}
	}
	json += "]"
	json += "}"

	return json
}

func (m *Requirement) asJSON() string {
	json := "{"
	json += fmt.Sprintf("\"name\":\"%v\",", m.Name)
	json += fmt.Sprintf("\"version_requirement\":\"%v\",", m.Version_requirement)

	json += "}"
	return json
}

func (m *Dependency) asJSON() string {
	json := "{"
	json += fmt.Sprintf("\"name\":\"%v\",", m.Name)
	json += fmt.Sprintf("\"version_requirement\":\"%v\",", m.Version_requirement)

	json += "}"
	return json
}

func (m *Metadata) asJSON() string {
	json := "{"
	json += fmt.Sprintf("\"name\":\"%v\",", m.Name)
	json += fmt.Sprintf("\"version\":\"%v\",", m.Version)
	json += fmt.Sprintf("\"author\":\"%v\",", m.Author)
	json += fmt.Sprintf("\"summary\":\"%v\",", m.Summary)
	json += fmt.Sprintf("\"license\":\"%v\",", m.License)
	json += fmt.Sprintf("\"source\":\"%v\",", m.Source)
	json += fmt.Sprintf("\"project_page\":\"%v\",", m.Project_page)

	if len(m.Issues_url) == 0 {
		json += "\"issues_url\":null,"
	} else {
		json += fmt.Sprintf("\"issues_url\":\"%v\",", m.Issues_url)
	}
	json += "\"dependencies\":["
	dependencies := m.Dependencies
	size := len(dependencies)
	for index, dependency := range dependencies {
		dependency_json := dependency.asJSON()
		json += dependency_json
		if index < size-1 {
			json += ","
		}
	}
	json += "]"

	if len(m.Data_provider) == 0 {
		json += "\"data_provider\":null,"
	} else {
		json += fmt.Sprintf("\"data_perovider\":\"%v\",", m.Data_provider)
	}

	json += "\"operatingsystem_support\":["
	operatingsystem_supports := m.Operatingsystem_support
	size = len(operatingsystem_supports)
	for index, operatingsystem_support := range operatingsystem_supports {
		operatingsystem_support_json := operatingsystem_support.asJSON()
		json += operatingsystem_support_json
		if index < size-1 {
			json += ","
		}
	}
	json += "]"

	json += "\"requirements\":["
	requirements := m.Requirements
	size = len(requirements)
	for index, requirement := range requirements {
		requirement_json := requirement.asJSON()
		json += requirement_json
		if index < size-1 {
			json += ","
		}
	}
	json += "]"

	json += "}"
	return json
}
