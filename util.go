package main

import (
	"fmt"
	"regexp"
)

func ValidModuleReleaseFilename(moduleReleaseFilename string) (bool, error) {

	result, _ := regexp.MatchString("^[a-zA-Z0-9]+[-/][a-z][a-z0-9_]*[-/][0-9]+.[0-9]+.[0-9]+(?:[-+].+)?.tar.gz$", moduleReleaseFilename)

	if result {
		return result, nil
	}
	return false, fmt.Errorf(`{"message":"400 Bad Request","errors":["'%s' is not a valid release slug"]}`, moduleReleaseFilename)
}

func ValidModuleReleaseSlug(release_slug string) (bool, error) {

	result, _ := regexp.MatchString("^[a-zA-Z0-9]+[-/][a-z][a-z0-9_]*[-/][0-9]+.[0-9]+.[0-9]+(?:[-+].+)?$", release_slug)

	if result {
		return result, nil
	}

	return false, fmt.Errorf(`{"message":"400 Bad Request","errors":["'%s' is not a valid release slug"]}`, release_slug)
}
