package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const DefaultCache string = "cache"

func FileInCache(filename string) bool {
	filePath := "cache/" + filename
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
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
