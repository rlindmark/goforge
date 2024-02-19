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

// func (m *OperatingSystemSupport) asJSON() string {
// 	json := "{"
// 	json += fmt.Sprintf(`%q:%q`, "operatingsystem", m.Operatingsystem)

// 	operatingsystemreleases := m.Operatingsystemrelease
// 	size := len(operatingsystemreleases)
// 	if size > 0 {
// 		json += `,"operatingsystemrelease":[`
// 		operatingsystemreleases := m.Operatingsystemrelease
// 		size := len(operatingsystemreleases)
// 		for index, operatingsystemrelease := range operatingsystemreleases {
// 			json += fmt.Sprintf(`%q`, operatingsystemrelease)
// 			if index < (size - 1) {
// 				json += ","
// 			}
// 		}
// 	}
// 	json += "]"
// 	json += "}"

// 	return json
// }

// func (m *Requirement) asJSON() string {
// 	json := "{"
// 	json += fmt.Sprintf(`%q:%q`, "name", m.Name)
// 	json += fmt.Sprintf(`%q:%q`, "version_requirement", m.Version_requirement)

// 	json += "}"
// 	return json
// }

// func (m *Dependency) asJSON() string {
// 	json := "{"
// 	json += fmt.Sprintf(`%q:%q`, "name", m.Name)
// 	json += fmt.Sprintf(`%q:%q`, "version_requirement", m.Version_requirement)

// 	json += "}"
// 	return json
// }

// func (m *Metadata) asJSON() string {
// 	json := "{"
// 	json += fmt.Sprintf(`%q:%q`, "name", m.Name)
// 	json += fmt.Sprintf(`%q:%q`, "version", m.Version)
// 	json += fmt.Sprintf(`%q:%q`, "author", m.Author)
// 	json += fmt.Sprintf(`%q:%q`, "summary", m.Summary)
// 	json += fmt.Sprintf(`%q:%q`, "license", m.License)
// 	json += fmt.Sprintf(`%q:%q`, "source", m.Source)
// 	json += fmt.Sprintf(`%q:%q`, "project_page", m.Project_page)

// 	if len(m.Issues_url) == 0 {
// 		json += `"issues_url":null,`
// 	} else {
// 		json += fmt.Sprintf(`%q:%q`, "issues_url", m.Issues_url)
// 	}
// 	json += `"dependencies":[`
// 	dependencies := m.Dependencies
// 	size := len(dependencies)
// 	for index, dependency := range dependencies {
// 		dependency_json := dependency.asJSON()
// 		json += dependency_json
// 		if index < size-1 {
// 			json += ","
// 		}
// 	}
// 	json += "]"

// 	if len(m.Data_provider) == 0 {
// 		json += `"data_provider":null,`
// 	} else {
// 		json += fmt.Sprintf(`%q:%q`, "data_perovider", m.Data_provider)
// 	}

// 	json += `"operatingsystem_support":[`
// 	operatingsystem_supports := m.Operatingsystem_support
// 	size = len(operatingsystem_supports)
// 	for index, operatingsystem_support := range operatingsystem_supports {
// 		operatingsystem_support_json := operatingsystem_support.asJSON()
// 		json += operatingsystem_support_json
// 		if index < (size - 1) {
// 			json += ","
// 		}
// 	}
// 	json += "]"

// 	json += `"requirements":[`
// 	requirements := m.Requirements
// 	size = len(requirements)
// 	for index, requirement := range requirements {
// 		requirement_json := requirement.asJson()
// 		json += requirement_json
// 		if index < size-1 {
// 			json += ","
// 		}
// 	}
// 	json += "]"
// 	json += "}"
// 	return json
// }
