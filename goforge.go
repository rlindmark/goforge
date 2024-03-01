package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const DefaultForgeIp = "127.0.0.1"
const DefaultForgePort = "8080"
const DefaultCacheRoot = "cache"

const DefaultPageLimit = 20
const DefaultPageOffset = 0

type ForgeError struct {
	Error_msg string   `json:"error"`
	Messages  []string `json:"messages"`
}

type V3ReleaseResponse struct {
	Pagination *Pagination    `json:"pagination"`
	Results    []PuppetModule `json:"results"`
}

func validModuleReleaseFilename(moduleReleaseFilename string) (bool, error) {

	result, _ := regexp.MatchString("^[a-zA-Z0-9]+[-/][a-z][a-z0-9_]*[-/][0-9]+.[0-9]+.[0-9]+(?:[-+].+)?.tar.gz$", moduleReleaseFilename)

	if result {
		return result, nil
	}
	return false, fmt.Errorf(`{"message":"400 Bad Request","errors":["'%s' is not a valid release slug"]}`, moduleReleaseFilename)
}

func validModuleReleaseSlug(release_slug string) (bool, error) {

	result, _ := regexp.MatchString("^[a-zA-Z0-9]+[-/][a-z][a-z0-9_]*[-/][0-9]+.[0-9]+.[0-9]+(?:[-+].+)?$", release_slug)

	if result {
		return result, nil
	}

	return false, fmt.Errorf(`{"message":"400 Bad Request","errors":["'%s' is not a valid release slug"]}`, release_slug)
}

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
	res, err := validModuleReleaseFilename(moduleReleaseFilename)
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

// SplitModuleName takes a valid module name and returns owner, module, version and error
// a module name might end in tar.gz which are disregarded during split.
// Example:
//
//	given owner-module-1.0.0 or owner-module-1.0.0.tar.gz returns
//	(owner, module, 1.0.0, nil)
func SplitModuleName(puppetmodule string) (string, string, string, error) {

	// puppetmodule ending with ".tar.gz"
	module_name := strings.TrimSuffix(puppetmodule, ".tar.gz")
	ok, err := validModuleReleaseSlug(module_name)
	if ok {
		l := strings.SplitN(module_name, "-", 3)
		if len(l) == 3 {
			return l[0], l[1], l[2], nil
		}
	}
	return "", "", "", err
}

func to_owner_and_modulename(puppet_module_without_version string) (string, string) {
	z := strings.SplitN(puppet_module_without_version, "-", 2)
	if len(z) < 2 {
		return z[0], ""
	}
	return z[0], z[1]
}

// func get_all_versions_for_module(module_name string) []string {

// 	base := Forge_cache()

// 	if len(module_name) == 0 {
// 		return nil
// 	}
// 	dir_hash := Module_hash(module_name)
// 	owner, _ := to_owner_and_modulename(module_name)
// 	path := base + "/" + dir_hash + "/" + owner

// 	old, _ := os.Getwd()
// 	os.Chdir(path)
// 	files, _ := filepath.Glob(module_name + "*.tar.gz")
// 	os.Chdir(old)

// 	//fmt.Printf("files: %v\n", files)
// 	return files
// }

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
	//fmt.Printf("ListModuleReleases:json:\n%s\n", string(jSON))
	fmt.Fprint(w, string(jSON))
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

	res, err := validModuleReleaseSlug(moduleReleaseSlug)
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

func Forge_ip() string {
	forge_ip := os.Getenv("FORGE_IP")
	if len(forge_ip) == 0 {
		return DefaultForgeIp
	}
	return forge_ip
}

func Forge_port() string {
	forge_port := os.Getenv("FORGE_PORT")
	if len(forge_port) == 0 {
		return DefaultForgePort
	}
	return forge_port
}

func Forge_cache() string {
	forge_cache := os.Getenv("FORGE_CACHE")
	if len(forge_cache) == 0 {
		return DefaultCacheRoot
	}
	return forge_cache
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

var help = flag.Bool("help", false, `Usage: use environment variables to configure usage.
	FORGE_IP - ip-adress to bind to (default 127.0.0.1)
	FORGE_PORT - listen port (default 8080)
	FORGE_CACHE - path to the forge cache (default 'cache')`)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

var forge_cache ForgeCache

func main() {

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	forge_cache = NewForgeCache(Forge_cache())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HomePage(w, r)
	})

	http.HandleFunc("/v3/releases/", func(w http.ResponseWriter, r *http.Request) {
		FetchModuleRelease(w, r)
	})

	http.HandleFunc("/v3/releases", func(w http.ResponseWriter, r *http.Request) {
		ListModuleReleases(w, r)
	})

	http.HandleFunc("/v3/files/", func(w http.ResponseWriter, r *http.Request) {
		DownloadModuleRelease(w, r)
	})

	ip_and_port := fmt.Sprintf("%s:%s", Forge_ip(), Forge_port())

	log.Fatal(http.ListenAndServe(ip_and_port, logRequest(http.DefaultServeMux)))
}
