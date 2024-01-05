package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"regexp"
)

const DefaultPort string = "8080"
const DefaultIP string = "127.0.0.1"

func ForgePort() string {

	return DefaultPort
}

func ForgeIP() string {

	return DefaultIP
}

func validModuleReleaseFilename(moduleReleaseFilename string) (bool, error) {

	result, _ := regexp.MatchString("^[a-zA-Z0-9]+[-/][a-z][a-z0-9_]*[-/][0-9]+.[0-9]+.[0-9]+(?:[-+].+)?.tar.gz$", moduleReleaseFilename)

	if result {
		return result, nil
	}
	return false, fmt.Errorf("{\"message\": \"400 Bad Request\", \"errors\": [\"'%s' is not a valid release slug\"]}", moduleReleaseFilename)
}

func handleV3Files(w http.ResponseWriter, r *http.Request) {
	// get filename
	moduleReleaseFilename := r.URL.Path[10:]
	// FIXME: test that r.URL.Path[0-9] is "/v3/files/"

	// test that filename is legal
	// FIXME: hansle error from function also. Maybe return the result json as the error
	res, err := validModuleReleaseFilename(moduleReleaseFilename)
	if !res {
		// 400
		result := fmt.Sprintf("{\"message\": \"400 Bad Request\", \"errors\": [\"'%s' is not a valid release slug\"]}", moduleReleaseFilename)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, result)
		return
	}

	// check if filename exist
	path, err := ModuleReleaseFilenameInCache(moduleReleaseFilename)

	// get filename from cache or error
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		// 404
		result := "{\"message\": \"404 Not Found\", \"errors\": [\"The requested resource could not be found\"]}"
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, result)
		return
	}

	// set Contant-Type:
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/v3/releases", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/v3/releases")
	})

	http.HandleFunc("/v3/files/", func(w http.ResponseWriter, r *http.Request) {
		handleV3Files(w, r)
	})

	ip_and_port := fmt.Sprintf("%s:%s", ForgeIP(), ForgePort())

	log.Fatal(http.ListenAndServe(ip_and_port, nil))
}
