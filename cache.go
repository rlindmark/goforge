package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const DefaultCache string = "cache"

// FIXME: argument need to be a file
func FileInCache(filename string) bool {
	filePath := filename
	fileInfo, error := os.Stat(filePath)
	if error != nil {
		return !errors.Is(error, os.ErrNotExist)
	}

	return fileInfo.Mode().IsRegular()
}

func ModulePathInCache(module_name string) (string, error) {
	if len(module_name) < 1 {
		return module_name, fmt.Errorf("module_name to short:%v", module_name)
	}
	cache := module_name[0]
	module_owner := strings.Split(module_name, "-")[0]
	return fmt.Sprintf("%s/%c/%s/%s", DefaultCache, cache, module_owner, module_name), nil
}

// Module release filename to be downloaded (e.g. "puppetlabs-apache-2.0.0.tar.gz")
func ModuleReleaseFilenameInCache(module_name string) (string, error) {
	if len(module_name) < 1 {
		return module_name, fmt.Errorf("module_name to short:%v", module_name)
	}
	cache := module_name[0]
	module_owner := strings.Split(module_name, "-")[0]
	return fmt.Sprintf("%s/%c/%s/%s", DefaultCache, cache, module_owner, module_name), nil
}

// func FindModulesInCache()

// true if filename = "puppetlabs-stdlib-1.0.0.tar.gz" exists at "cache/p/puppetlabs/puppetlabs-stdlib-1.0.0.tar.gz" and is a file
// true if filename = "p/puppetlabs/puppetlabs-stdlib-1.0.0.tar.gz" exists at "cache/p/puppetlabs/puppetlabs-stdlib-1.0.0.tar.gz" and is a file
func InCache(filename string) bool {
	return false
}

func GetItemPath(filename string) (string, error) {
	return "", fmt.Errorf("%s does not exist", filename)
}

func StoredModules() (fileList []string) {
	fmt.Print("StoredModules: TO BE IMPLEMENTED")
	return []string{}
}
