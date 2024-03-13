package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const DefaultCache string = "cache"

type ForgeCache struct {
	cache_root string
}

func NewForgeCache(root_path string) ForgeCache {
	return ForgeCache{cache_root: root_path}
}

// FIXME: argument need to be a file
func FileInCache(filename string) bool {
	filePath := filename
	fileInfo, error := os.Stat(filePath)
	if error != nil {
		return !errors.Is(error, os.ErrNotExist)
	}

	return fileInfo.Mode().IsRegular()
}

func Module_hash(module string) string {
	if len(module) > 0 {
		return string(module[0])
	}
	return ""
}

// GetModule returns the filepath to the module if found in the cache
func (c *ForgeCache) GetModuleFilePath(release_slug string) (string, error) {

	ok, err := validModuleReleaseSlug(release_slug)
	if !ok {
		return "", err
	}

	owner, module, version, err := SplitModuleName(release_slug)
	if err != nil {
		return "", err
	}

	hash_path := Module_hash(release_slug)
	path := filepath.Join(c.cache_root, hash_path, owner, module, version)

	if FileInCache(path + ".tar.gz") {
		return path + ".tar.gz", nil
	}
	return "", fmt.Errorf("not found in cache")
}

// GetModuleVersions returns a list of filepaths for the module slug in the cache
func (c ForgeCache) GetModuleVersions(module_slug string) []string {

	//base := Forge_cache()

	ok, _ := isValidModuleSlug(module_slug)
	if !ok {
		return nil
	}

	hashpath := Module_hash(module_slug)
	owner, _ := to_owner_and_modulename(module_slug)
	path := filepath.Join(c.cache_root, hashpath, owner)

	old, _ := os.Getwd()
	os.Chdir(path)
	files, _ := filepath.Glob(module_slug + "*.tar.gz")
	os.Chdir(old)

	//fmt.Printf("files: %v\n", files)
	return files
}

// GetAllUsers returns a list of all users found in the cache
func (c ForgeCache) GetAllUsers() []string {

	//base := Forge_cache()

	files, err := filepath.Glob("cache/*/*")
	if err != nil {
		log.Fatal(err)
	}

	vsm := make([]string, len(files))
	for i, v := range files {
		vsm[i] = filepath.Base(v)
	}
	return vsm
	// FIXME: only return directories and not files
	// for _, file := range files {
	// 	fmt.Println(file)
	// }

	// return files
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
