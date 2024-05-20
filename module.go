package main

type Module struct {
	Uri            string  `json:"uri"`
	Slug           string  `json:"slug"`
	Name           string  `json:"name"`
	Downloads      int     `json:"downloads"`
	Created_at     string  `json:"created_at"`
	Updated_at     string  `json:"updated_at"`
	Deprecated_at  string  `json:"deprecated_at"`
	Deprecated_for *string `json:"deprecated_for"`
	Superseeded_by *string `json:"superseeded_by"`
	Supported      bool    `json:"supported"`
	Module_group   string  `json:"module_group"`
	Premium        bool    `json:"premium"`

	Owner           *Owner        `json:"owner"`
	Current_Release *PuppetModule `json:"current_release"`
	Releases        []Releases    `json:"releases"`
	Feedback_score  int           `json:"feedback_score"`
	Homepage_url    string        `json:"homepage_url"`
	Issues_url      string        `json:"issues_url"`
}

type Releases struct {
	Uri        string  `json:"uri"`
	Slug       string  `json:"slug"`
	Version    string  `json:"version"`
	Supported  bool    `json:"supported"`
	Created_at string  `json:"created_at"`
	Deleted_at *string `json:"deleted_at"`
	File_uri   string  `json:"file_uri"`
	File_size  int     `json:"file_size"`
}
