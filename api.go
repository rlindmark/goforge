package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
DownloadModuleRelease downloads filename

GET /v3/files/{filename}

# PATH PARAMETERS

	filename (required) Module release filename to be downloaded (e.g. "puppetlabs-apache-2.0.0.tar.gz")
*/
func DownloadModuleRelease(w http.ResponseWriter, r *http.Request) {

	// moduleReleaseFilename are on the form puppetlabs-apache-4.0.0
	// Module release filename to be downloaded (e.g. "puppetlabs-apache-2.0.0.tar.gz")
	moduleReleaseFilename := r.URL.Path[10:]
	// FIXME: test that r.URL.Path[0-9] is "/v3/files/"

	// test that filename is legal
	res, err := ValidModuleReleaseFilename(moduleReleaseFilename)
	if !res {
		//result := fmt.Sprintf(`{"message":"400 Bad Request","errors":["'%s' is not a valid release slug"]}`, moduleReleaseFilename)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	// check if filename exist
	path, _ := ModuleReleaseFilenameInCache(moduleReleaseFilename)

	// get filename from cache or error
	file, err := os.Open(path)

	if err != nil {
		// 404
		result := `{"message":"404 Not Found","errors":["The requested resource could not be found"]}`
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, result)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream; charset=utf-8")
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

/*
ListModuleReleases returns a list of module releases meeting the specified search criteria and filters. Results are paginated. All of the parameters are optional.

GET /v3/releases

QUERY PARAMETERS

	limit  integer [1..100} Default: 20

	offset integer >= 0 Default: 0

	sort_by Enum["downloads" "release_date" "module"] Desired order in which to return results

	module
	owner  NOT IMPLEMENTED
	...    NOT IMPLEMENTED
*/
func ListModuleReleases(w http.ResponseWriter, r *http.Request) {

	url_query := r.URL.Query()

	// parse all query parameters
	// FIXME: there are a lot to manage
	module_name := r.URL.Query().Get("module")

	var err error

	offset := DefaultPageOffset
	offset_string := r.URL.Query().Get("offset")
	if offset_string != "" {
		offset, err = strconv.Atoi(offset_string)
		if err != nil {
			// if not an integer, report it and let offset = defaultPageOffset
			// FIXME: this should return an error instead
			fmt.Printf("expected integer, got %v", offset_string)
		}
	} else {
		// if offset is not present in the query string, add it with default value
		url_query.Add("offset", fmt.Sprint(offset))
	}

	limit := DefaultPageLimit
	limit_string := r.URL.Query().Get("limit")
	if limit_string != "" {
		// FIXME: check for err
		limit, err = strconv.Atoi(limit_string)
		if err != nil {
			// if not an integer, report it and let offset = defaultPageOffset
			// FIXME: this should return an error instead
			fmt.Printf("expected integer, got %v", offset_string)
		}
	} else {
		// if limit is not present in the query string, add it with default value
		url_query.Add("limit", fmt.Sprint(limit))
	}

	//fmt.Printf("module_name: %v\n", module_name)

	// all_matching_modules := get_all_versions_for_module(module_name)
	// fmt.Printf(forge_cache.cache_root)
	all_matching_modules := forge_cache.GetModuleVersions(module_name)

	sort.Sort(sort.Reverse(sort.StringSlice(all_matching_modules)))
	total := len(all_matching_modules)

	pagination, _ := CreatePagination(url_query, total)

	modules_list, _ := get_results(all_matching_modules, offset, limit)

	response := V3ReleaseResponse{Pagination: pagination, Results: modules_list}

	jSON, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Unable to marshal. error:%v", err)
		return
	}
	fmt.Fprint(w, string(jSON))
}

// CreateModuleRelease publish a new module or new release of an existing module
func CreateModuleRelease(w http.ResponseWriter, r *http.Request) {

}

/*
FetchModuleRelease returns data for a single module Release resource identified by the module release's slug value.

GET /v3/releases/{release_slug}

# PATH PARAMETERS

	release_slug (required) example: puppetlabs-apache-4.0.0

QUERY PARAMETERS

	with_html        NOT IMPLEMENTED
	include_fields   NOT IMPLEMENTED
	exclude_fields   NOT IMPLEMENTED
*/
func FetchModuleRelease(w http.ResponseWriter, r *http.Request) {

	// FIXME: ensure the first 13 bytes in r.URL.Path
	moduleReleaseSlug := r.URL.Path[13:]

	res, err := ValidModuleReleaseSlug(moduleReleaseSlug)
	if !res {
		// 400
		//result := fmt.Sprintf(`{"message":"400 Bad Request","errors":["'%s' is not a valid release slug"]}`, moduleReleaseSlug)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	module, err := NewPuppetModule(moduleReleaseSlug)
	if module == nil {
		// 404
		//result := `{"message":"404 Not Found","errors":["'The requested resource could not be found"]}`
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
		return
	}
	jSON, _ := json.Marshal(module)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//fmt.Printf("FetchModuleRelease:json:%v", string(jSON))
	fmt.Fprint(w, string(jSON))
}

func DeleteModuleRelease(w http.ResponseWriter, r *http.Request) {

}

func ListModuleReleasePlans(w http.ResponseWriter, r *http.Request) {

}

func FetchModuleReleasePlan(w http.ResponseWriter, r *http.Request) {

}

func ListModules(w http.ResponseWriter, r *http.Request) {

}

func FetchModule(w http.ResponseWriter, r *http.Request) {

}

func DeprecateModule(w http.ResponseWriter, r *http.Request) {

}

func DeleteModule(w http.ResponseWriter, r *http.Request) {

}

func CreateSearchFilter(w http.ResponseWriter, r *http.Request) {

}

func GetUsersSearchFilters(w http.ResponseWriter, r *http.Request) {

}

func DeleteSearchFilterByID(w http.ResponseWriter, r *http.Request) {

}

/*
ListUsers provides information about Puppet Forge user accounts. By default,
results are returned in alphabetical order by username and paginated with 20
users per page. It's also possible to sort by number of published releases,
total download counts for all the user's modules, or by the date of the
user's latest release. All parameters are optional.
*/
func ListUsers(w http.ResponseWriter, r *http.Request) {

	type userResponse struct {
		Pagination *Pagination `json:"pagination"`
		Results    []User      `json:"results"`
	}

	url_query := r.URL.Query()

	// parse all query parameters
	// FIXME: there are a lot to manage

	var err error

	offset := DefaultPageOffset
	offset_string := r.URL.Query().Get("offset")
	if offset_string != "" {
		offset, err = strconv.Atoi(offset_string)
		if err != nil {
			// if not an integer, report it and let offset = defaultPageOffset
			// FIXME: this should return an error instead
			fmt.Printf("expected integer, got %v", offset_string)
		}
	} else {
		// if offset is not present in the query string, add it with default value
		url_query.Add("offset", fmt.Sprint(offset))
	}

	limit := DefaultPageLimit
	limit_string := r.URL.Query().Get("limit")
	if limit_string != "" {
		// FIXME: check for err
		limit, err = strconv.Atoi(limit_string)
		if err != nil {
			// if not an integer, report it and let offset = defaultPageOffset
			// FIXME: this should return an error instead
			fmt.Printf("expected integer, got %v", offset_string)
		}
	} else {
		// if limit is not present in the query string, add it with default value
		url_query.Add("limit", fmt.Sprint(limit))
	}

	users := forge_cache.GetAllUsers()

	sort.Sort(sort.Reverse(sort.StringSlice(users)))
	total := len(users)

	pagination, _ := CreatePagination(url_query, total)

	user_list, _ := get_user_results(users, offset, limit)

	response := userResponse{Pagination: pagination, Results: user_list}

	jSON, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Unable to marshal. error:%v", err)
		return
	}
	fmt.Fprint(w, string(jSON))
}

/*
FetchUser returns data for a single User resource identified by the user's slug value.

GET /v3/users/{user_slug}
PATH PARAMETERS

user_slug required string  ^[a-zA-Z0-9]+$ example: puppetlabs

Unique textual identifier (slug) of the User resource to retrieve
*/
func FetchUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		return
	}

	userSlug := r.URL.Path[10:]

	res, _ := IsValidUserSlug(userSlug)
	if !res {
		// 400
		result := fmt.Sprintf(`{"message":"400 Bad Request","errors":["'%s' is not a valid user slug"]}`, userSlug)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, result)
		return
	}

	user, err := NewUser("/v3/users/"+userSlug, userSlug, "12345", userSlug, userSlug, 1, 1, "1970-01-01 01:01", "1970-01-01 01:01")
	if user == nil {
		// 404
		// result := `{"message":"404 Not Found","errors":["'The requested resource could not be found"]}`
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
		return
	}

	jSON, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	fmt.Fprint(w, string(jSON))
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func get_user_results(users []string, offset int, limit int) ([]User, error) {

	// assert first >= 0
	// assert last >= first
	// len(all_modules) >= last

	var result []User
	total := len(users)

	if total == 0 {
		return nil, nil
	}
	last := min(total, offset+limit)

	create_at := "1970-01-01 01:01:01 0000"  // just make some up
	updated_at := "1970-01-01 01:01:01 0000" // just make some up
	gravatar_id := "1234"

	for _, user_name := range users[offset:last] {
		user, err := NewUser("/v3/user/"+user_name, user_name, gravatar_id, user_name, user_name, 0, 0, create_at, updated_at)
		if err == nil {
			result = append(result, *user)
		}
	}
	return result, nil
}
func get_results(all_modules []string, offset int, limit int) ([]PuppetModule, error) {

	// assert first >= 0
	// assert last >= first
	// len(all_modules) >= last

	var result []PuppetModule
	total := len(all_modules)

	if total == 0 {
		return nil, nil
	}
	last := min(total, offset+limit)

	for _, module_name_version := range all_modules[offset:last] {
		puppet_module, _ := get_v3_releases_module_result(module_name_version)
		result = append(result, *puppet_module)
	}
	return result, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func get_v3_releases_module_result(module_name string) (*PuppetModule, error) {

	modulename := strings.TrimSuffix(module_name, ".tar.gz")

	puppet_module, _ := NewPuppetModule(modulename)

	return puppet_module, nil
}
