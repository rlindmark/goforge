package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

/*
DownloadModuleRelease downloads filename

GET /v3/files/{filename}

PATH PARAMETERS

	filename (required) Module release filename to be downloaded (e.g. "puppetlabs-apache-2.0.0.tar.gz")
*/
func DownloadModuleRelease(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		return
	}

	if filepath.Dir(r.URL.Path) != "/v3/files" {
		result := fmt.Sprintf(`{"message":"500 Internal Server Error","errors":["Internal Server Error. Path=%v"]}`, r.URL.Path)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result)
		return
	}

	moduleReleaseFilename := r.URL.Path[len("/v3/files/"):]

	res, err := ValidModuleReleaseFilename(moduleReleaseFilename)
	if !res {
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

// HandleReleases manages the /v3/releases endpoint
func HandleReleases(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ListModuleReleases(w, r)
	case "POST":
		CreateModuleRelease(w, r)
	default:
		// NOT implemented
	}
}

/*
ListModuleReleases returns a list of module releases meeting the specified search criteria and filters.
Results are paginated. All of the parameters are optional.

GET /v3/releases

QUERY PARAMETERS

	limit  integer [1..100] Default: 20

	offset integer >= 0 Default: 0

	sort_by Enum["downloads" "release_date" "module"] Desired order in which to return results NOT IMPLEMENTED

	module
	owner  NOT IMPLEMENTED
	...    NOT IMPLEMENTED
*/
func ListModuleReleases(w http.ResponseWriter, r *http.Request) {

	type releaseResponse struct {
		Pagination *Pagination    `json:"pagination"`
		Results    []PuppetModule `json:"results"`
	}

	// only handle GET requests
	if r.Method != "GET" {
		return
	}

	url_query := r.URL.Query()
	path := r.URL.Path

	// parse all query parameters
	// FIXME: there are a lot to manage
	module_name := r.URL.Query().Get("module")

	var err error
	var offset int
	var limit int

	offset_string := r.URL.Query().Get("offset")
	if offset_string == "" {
		// offset is not present in the query string, add it with default value
		offset = DefaultPageOffset
		url_query.Add("offset", fmt.Sprint(offset))
	} else {
		offset, err = strconv.Atoi(offset_string)
		if err != nil {
			// 400 BadRequest
			result := fmt.Sprintf(`{"message":"400 Bad Request","errors":["offset '%s' is not an integer"]}`, offset_string)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, result)
			return
		}
	}

	limit_string := r.URL.Query().Get("limit")
	if limit_string == "" {
		// if limit is not present in the query string, add it with default value
		limit = DefaultPageLimit
		url_query.Add("limit", fmt.Sprint(limit))
	} else {
		limit, err = strconv.Atoi(limit_string)
		if err != nil {
			// 400 BadRequest
			result := fmt.Sprintf(`{"message":"400 Bad Request","errors":["limit '%s' is not an integer"]}`, limit_string)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, result)
			return
		}
	}

	all_matching_modules := forge_cache.GetModuleVersions(module_name)

	sort.Sort(sort.Reverse(sort.StringSlice(all_matching_modules)))
	total := len(all_matching_modules)

	pagination, err := CreatePagination(path, url_query, total)
	if err != nil {
		result := fmt.Sprintf(`{"message":"400 Bad Request","errors":["%v"]}`, err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, result)
		return
	}

	modules, _ := get_PuppetModules(all_matching_modules, offset, limit)

	response := releaseResponse{Pagination: pagination, Results: modules}

	jSON, err := json.Marshal(response)
	if err != nil {
		// FIXME: handle Marshal error
		fmt.Printf("Unable to marshal. error:%v", err)
		return
	}
	fmt.Fprint(w, string(jSON))
}

// CreateModuleRelease publish a new module or new release of an existing module
// NOTE: This function is only a placeholder
func CreateModuleRelease(w http.ResponseWriter, r *http.Request) {

	// only handle POST requests
	if r.Method != "POST" {
		return
	}

	// Check for "Authorization: Bearer <api_key>"" header
	authorization := r.Header.Get("Authorization")
	if authorization == "" {

		result := `{"message":"401 Unauthorized","errors":["This endpoint requires a valid Authorization header"]}`
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, result)
		return
	}

	// Return 403 as no upload are allowed
	result := `{"message":"403 Forbidden","errors":["The provided API key is invalid or has insufficient permissions for the requested operation"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, result)
}

// HandleModuleRelease manages the /v3/releases/{release_slug} endpoint
func HandleModuleRelease(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		FetchModuleRelease(w, r)
	case "DELETE":
		DeleteModuleRelease(w, r)
	default:
		// NOT implemented
	}
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

	// only handle GET requests
	if r.Method != "GET" {
		return
	}

	if filepath.Dir(r.URL.Path) != "/v3/releases" {
		result := fmt.Sprintf(`{"message":"500 Internal Server Error","errors":["Internal Server Error. Path=%v"]}`, r.URL.Path)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result)
		return
	}

	moduleReleaseSlug := r.URL.Path[len("/v3/releases/"):]

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

	// only handle DELETE requests
	if r.Method != "DELETE" {
		return
	}

	// Check for "Authorization: Bearer <api_key>"" header
	authorization := r.Header.Get("Authorization")
	if authorization == "" {

		result := `{"message":"401 Unauthorized","errors":["This endpoint requires a valid Authorization header"]}`
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, result)
		return
	}

	// Return 403 as no delete are allowed
	result := `{"message":"403 Forbidden","errors":["The provided API key is invalid or has insufficient permissions for the requested operation"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, result)
}

func ListModuleReleasePlans(w http.ResponseWriter, r *http.Request) {

	// only handle GET requests
	if r.Method != "GET" {
		return
	}

	// NOT IMPLEMENTED
	result := `{"message":"404 Not Found","errors":["The requested resource could not be found"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, result)
}

func FetchModuleReleasePlan(w http.ResponseWriter, r *http.Request) {

	// only handle GET requests
	if r.Method != "GET" {
		return
	}
	// NOT IMPLEMENTED
	result := `{"message":"404 Not Found","errors":["The requested resource could not be found"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, result)
}

/*
ListModules returns a list of modules meeting the specified search criteria and filters.
Results are paginated. All of the parameters are optional.
To publish or delete a Release resource, see Release operations.
*/
func ListModules(w http.ResponseWriter, r *http.Request) {

	// only handle GET requests
	if r.Method != "GET" {
		return
	}

	// NOT IMPLEMENTED
	result := `{"message":"404 Not Found","errors":["The requested resource could not be found"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, result)
}

func IsValidModuleSlug(module_slug string) (bool, error) {
	return regexp.MatchString("^[a-zA-Z0-9]+[-/][a-z][a-z0-9_]*$", module_slug)
}

func HandleModules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		FetchModule(w, r)
	case "PATCH":
		DeprecateModule(w, r)
	case "DELETE":
		DeleteModule(w, r)
	default:
		// NOT implemented
	}
}

/*
FetchModule returns data for a single Module resource identified by the module's slug value.

GET /v3/modules/{module_slug}

PATH PARAMETERS
module_slug required string ^[a-zA-Z0-9]+[-\/][a-z][a-z0-9_]*$
Example: puppetlabs-apache

QUERY PARAMETERS
with_html boolean Render markdown files (README, REFERENCE, etc.) to HTML before returning results

include_fields Array of strings
Example: include_fields=docs
List of top level keys to include in response object, only applies to fields marked 'optional'

exclude_fields Array of strings
Example: exclude_fields=readme changelog license reference
List of top level keys to exclude from response object
*/
func FetchModule(w http.ResponseWriter, r *http.Request) {

	// only handle GET requests
	if r.Method != "GET" {
		return
	}

	module_slug := r.URL.Path[10:]

	res, _ := IsValidModuleSlug(module_slug)
	if !res {
		// 400
		result := fmt.Sprintf(`{"message":"400 Bad Request","errors":["'%s' is not a valid release slug"]}`, module_slug)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, result)
		return
	}

	result := `{"message":"404 Not Found","errors":["The requested resource could not be found"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, result)
}

func DeprecateModule(w http.ResponseWriter, r *http.Request) {

	// only handle PATCH requests
	if r.Method != "PATCH" {
		return
	}
	// Check for "Authorization: Bearer <api_key>"" header
	authorization := r.Header.Get("Authorization")
	if authorization == "" {

		result := `{"message":"401 Unauthorized","errors":["This endpoint requires a valid Authorization header"]}`
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, result)
		return
	}

	// Return 403 as no delete are allowed
	result := `{"message":"403 Forbidden","errors":["The provided API key is invalid or has insufficient permissions for the requested operation"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, result)
}

func DeleteModule(w http.ResponseWriter, r *http.Request) {

	// only handle DELETE requests
	if r.Method != "DELETE" {
		return
	}
	// Check for "Authorization: Bearer <api_key>"" header
	authorization := r.Header.Get("Authorization")
	if authorization == "" {

		result := `{"message":"401 Unauthorized","errors":["This endpoint requires a valid Authorization header"]}`
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, result)
		return
	}

	// Return 403 as no delete are allowed
	result := `{"message":"403 Forbidden","errors":["The provided API key is invalid or has insufficient permissions for the requested operation"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, result)

}

func HandleSearchFilters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		CreateSearchFilter(w, r)
	case "GET":
		GetUsersSearchFilters(w, r)
	default:
		// NOT implemented
	}
}

func CreateSearchFilter(w http.ResponseWriter, r *http.Request) {

	// only handle POST requests
	if r.Method != "POST" {
		return
	}

	result := `{"message":"404 Not Found","errors":["The requested resource could not be found"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, result)
}

func GetUsersSearchFilters(w http.ResponseWriter, r *http.Request) {

	// only handle POST requests
	if r.Method != "POST" {
		return
	}

	result := `{"message":"404 Not Found","errors":["The requested resource could not be found"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, result)
}

func DeleteSearchFilterByID(w http.ResponseWriter, r *http.Request) {

	// only handle DELETE requests
	if r.Method != "DELETE" {
		return
	}

	result := `{"message":"403 Forbidden","errors":["The provided API key is invalid or has insufficient permissions for the requested operation"]}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, result)
}

/*
ListUsers provides information about Puppet Forge user accounts.

By default, results are returned in alphabetical order by username and
paginated with 20 users per page. It's also possible to sort by number
of published releases, total download counts for all the user's modules,
or by the date of the user's latest release.

All parameters are optional.
*/
func ListUsers(w http.ResponseWriter, r *http.Request) {

	type userResponse struct {
		Pagination *Pagination `json:"pagination"`
		Results    []User      `json:"results"`
	}

	url_query := r.URL.Query()
	path := r.URL.Path
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
	if limit_string == "" {
		// if limit is not present in the query string, add it with default value
		url_query.Add("limit", fmt.Sprint(limit))
	} else {
		limit, err = strconv.Atoi(limit_string)
		if err != nil {
			// if not an integer, report it and let offset = defaultPageOffset
			fmt.Printf("expected integer, got %v", limit_string)
			url_query.Add("limit", fmt.Sprint(DefaultPageLimit))
		}

	}

	users := forge_cache.GetAllUsers()

	sort.Sort(sort.Reverse(sort.StringSlice(users)))
	total := len(users)

	pagination, err := CreatePagination(path, url_query, total)
	if err != nil {
		result := fmt.Sprintf(`{"message":"400 Bad Request","errors":["%v"]}`, err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, result)
		return
	}

	user_list, _ := get_user_results(users, offset, limit)

	response := userResponse{Pagination: pagination, Results: user_list}

	jSON, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Unable to marshal. error:%v", err)
		return
	}
	fmt.Fprint(w, string(jSON))
}

func GetGravatarId(username string) string {
	// FIXME: gravatar_id should be sha256 char long.

	return "nogravatardefined"
}

// GetModuleReleaseCount returns total number of released modules username have in cache
func GetModuleReleaseCount(username string) int {
	// FIXME: Implement code
	return 0
}

// GetModuleCount returns total number of modules username have in cache
func GetModuleCount(username string) int {
	// FIXME: Implement code
	return 0
}

/*
FetchUser returns data for a single User resource identified by the user's slug value.

GET /v3/users/{user_slug}
PATH PARAMETERS

user_slug required string ^[a-zA-Z0-9]+$ example: puppetlabs

# Unique textual identifier (slug) of the User resource to retrieve

# QUERY PARAMETERS

with_html	boolean Render markdown files (README, REFERENCE, etc.) to HTML before returning results

include_fields	Array of strings
Example: include_fields=docs
List of top level keys to include in response object, only applies to fields marked 'optional'

exclude_fields	Array of strings

Example: exclude_fields=readme changelog license reference
List of top level keys to exclude from response object
*/
func FetchUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		return
	}

	if filepath.Dir(r.URL.Path) != "/v3/users" {
		// This is a programming error. The calling code should make sure this does not happen
		result := fmt.Sprintf(`{"message":"500 Internal Server Error","errors":["Internal Server Error. Path=%v"]}`, r.URL.Path)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result)
		return
	}

	userSlug := r.URL.Path[len("/v3/users/"):]

	res, _ := IsValidUserSlug(userSlug)
	if !res {
		// 400
		result := fmt.Sprintf(`{"message":"400 Bad Request","errors":["'%s' is not a valid user slug"]}`, userSlug)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, result)
		return
	}

	// parse extra query parameters
	q := r.URL.Query()
	with_http := q.Get("with_http")
	fmt.Printf("with_http:%v", with_http)
	include_fields := q.Get("include_fields")
	fmt.Printf("include_fields:%v", include_fields)
	exclude_fields := q.Get("exclude_fields")
	fmt.Printf("exclude_fields:%v", exclude_fields)

	// FIXME: Need to check that the user exist in cache
	uri := "/v3/users/" + userSlug
	release_count := GetModuleReleaseCount(userSlug)
	module_count := GetModuleCount(userSlug)

	gravatar_id := GetGravatarId(userSlug)

	created_at := GetUserCreatedAt(userSlug)
	updated_at := GetUserUpdatedAt(userSlug)

	user, err := NewUser(uri, userSlug, gravatar_id, userSlug, userSlug, release_count, module_count, created_at, updated_at)
	if user == nil {
		// 404
		// result := `{"message":"404 Not Found","errors":["'The requested resource could not be found"]}`
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
		return
	}

	jSON, _ := json.Marshal(user)
	// FIXME: better to catch Marshal error here
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	fmt.Fprint(w, string(jSON))
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to goforge. A limited implementation of the puppet forgeapi, %q", html.EscapeString(r.URL.Path))
}

func GetUserCreatedAt(username string) string {
	return "1970-01-01 01:01:01 0000" // just make some up
}

func GetUserUpdatedAt(username string) string {
	return "1970-01-01 01:01:01 0000" // just make some up
}

// GetUserModuleCount returns total number of releases username have in cache
func GetUserReleaseCount(username string) int {
	// FIXME: not implemented. Duplicate of GetModuleReleaseCount()
	return 0
}

// GetUserModuleCount returns number of modules username have in cache
func GetUserModuleCount(username string) int {
	// FIXME: not implemented
	return 0
}

func get_user_results(users []string, offset int, limit int) ([]User, error) {

	// assert offset >= 0
	// assert last >= first
	// len(all_modules) >= last

	var result []User
	total := len(users)

	if total == 0 {
		return nil, nil
	}
	last := min(total, offset+limit)

	//release_count := 1
	//module_count := 1

	for _, username := range users[offset:last] {
		// FIXME: need to check that the user exist in cache

		uri := "/v3/user/" + username
		gravatar_id := GetGravatarId(username)
		release_count := GetModuleReleaseCount(username)
		module_count := GetModuleCount(username)
		created_at := GetUserCreatedAt(username)
		updated_at := GetUserUpdatedAt(username)

		user, err := NewUser(uri, username, gravatar_id, username, username, release_count, module_count, created_at, updated_at)
		if err == nil {
			result = append(result, *user)
		}
	}
	return result, nil
}

func get_PuppetModules(all_modules []string, offset int, limit int) ([]PuppetModule, error) {

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
