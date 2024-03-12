package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
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
