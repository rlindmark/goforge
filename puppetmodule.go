package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

type PuppetModule struct {
	Uri              string  `json:"uri"`  // "/v3/releases/puppetlabs-apache-4.0.0"
	Slug             string  `json:"slug"` // "puppetlabs-apache-4.0.0"
	Module           *Module `json:"module"`
	Version          string  `json:"version"`
	Metadata         string  `json:"metadata"`
	Tags             string  `json:"tags"`
	Supported        string  `json:"supported"`
	Pdk              string  `json:"pdk"`
	Validation_Score int     `json:"validation_score"`

	File_uri         string `json:"file_uri"` // "/v3/files/puppetlabs-apache-4.0.0.tar.gz"
	File_size        int64  `json:"file_size"`
	File_md5         string `json:"file_md5"`
	File_sha256      string `json:"file_sha256"`
	Downloads        int    `json:"downloads"`
	Readme           string `json:"readme"`
	Changelog        string `json:"changelog"`
	License          string `json:"license"`
	Reference        string `json:"reference"`
	Docs             string `json:"docs"`
	Pe_compatibility string `json:"pe_compatibility"`
	Tasks            string `json:"tasks"`
	Plans            string `json:"plans"`
}

func file_sha256(filename string) (string, error) {

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hash := sha256.New()
	defer f.Close()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func file_md5(filename string) (string, error) {

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hash := md5.New()
	defer f.Close()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func get_metadata(owner_module_version string) (string, error) {

	srcFile, _ := ModuleReleaseFilenameInCache(owner_module_version + ".tar.gz")
	//num := 1
	f, err := os.Open(srcFile)
	if err != nil {
		//fmt.Println(err)
		return "", err
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		//fmt.Println(err)
		return "", err
	}

	tarReader := tar.NewReader(gzf)

	i := 0
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			//fmt.Println(err)
			return "", err
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			//fmt.Println("(", i, ")", "Name: ", name)
			if name == owner_module_version+"/metadata.json" {
				metadata, err := io.ReadAll(tarReader)
				return string(metadata), err
			}
			// if i == num {
			// 	fmt.Println(" --- ")
			// 	io.Copy(os.Stdout, tarReader)
			// 	fmt.Println(" --- ")
			// 	os.Exit(0)
			// }
		default:
			fmt.Printf("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}

		i++
	}
	return "", nil
}

// func read_metadata(filename string) string {
// 	dat, err := os.ReadFile(filename)
// 	if err != nil {
// 		fmt.Printf("Cant open %s\n", filename)
// 		return ""
// 	}
// 	return (string(dat))
// }

// module_name is the string "owner-module-version" with no version
func NewPuppetModule(owner_module_version string) (*PuppetModule, error) {

	owner_module_version_tgz := fmt.Sprintf("%s.tar.gz", owner_module_version)

	tmp := strings.Split(owner_module_version, "-")
	module_owner := tmp[0]
	module_name := tmp[1]
	module_version := tmp[2]

	file_in_cache, _ := ModuleReleaseFilenameInCache(owner_module_version_tgz)

	uri := fmt.Sprintf("/v3/releases/%s", owner_module_version)
	//ok, _ := isValidModuleName(module_name)
	//filename := fmt.Sprintf("%s.tar.gz", owner_module_version)
	if !FileInCache(file_in_cache) {
		return nil, fmt.Errorf(`{"message":"404 Not Found","errors":["'The requested resource could not be found"]}`)
	}
	slug := owner_module_version
	module, _ := NewModule(fmt.Sprintf("%s-%s", module_owner, module_name))
	version := module_version
	//metadata := read_metadata(owner_module_version + "/metadata.json") //"{ }" // need to get this from metadata.json
	metadata, _ := get_metadata(owner_module_version) //"{ }" // need to get this from metadata.json
	tags := "[ string ]"
	supported := "true"
	pdk := "true"
	validation_score := 100 // dont know how to set this

	file_uri := fmt.Sprintf("/v3/files/%s", owner_module_version_tgz)

	fileInfo, err := os.Stat(file_in_cache)
	if err != nil {
		return nil, fmt.Errorf("could not os.stat(%s)", file_in_cache)
	}

	file_size := fileInfo.Size()
	file_md5, _ := file_md5(file_in_cache)
	file_sha256, _ := file_sha256(file_in_cache)
	downloads := 0
	readme := "Readme"
	changelog := "Changelog"
	license := "License"
	reference := "Reference"
	docs := "{ }"
	pe_compatibility := `[ "string" ]`
	tasks := "[ ]"
	plans := "[ ]"

	// if !ok {
	// 	return nil, fmt.Errorf("cant create module with invalid uri: %v", uri)
	// }
	puppetmodule := PuppetModule{uri, slug, module, version, metadata, tags, supported, pdk, validation_score, file_uri, file_size, file_md5,
		file_sha256, downloads, readme, changelog, license, reference, docs, pe_compatibility, tasks, plans}

	return &puppetmodule, nil
}

func (module *PuppetModule) asJson() string {
	json := "{"
	json += fmt.Sprintf("%q:%q,", "uri", module.Uri)
	json += fmt.Sprintf("%q:%q,", "slug", module.Slug)
	json += fmt.Sprintf("%q:%s,", "module", module.Module.asJson())
	json += fmt.Sprintf("%q:%q,", "version", module.Version)
	json += fmt.Sprintf("%q:%s,", "metadata", module.Metadata)
	json += fmt.Sprintf("%q:%q,", "tags", module.Tags)
	json += fmt.Sprintf("%q:%s,", "supported", module.Supported)
	json += fmt.Sprintf("%q:%s,", "pdk", module.Pdk)
	json += fmt.Sprintf("%q:%d,", "validation_score", module.Validation_Score)

	json += fmt.Sprintf("%q:%q,", "file_uri", module.File_uri)
	json += fmt.Sprintf("%q:%d,", "file_size", module.File_size)
	json += fmt.Sprintf("%q:%q,", "file_md5", module.File_md5)
	json += fmt.Sprintf("%q:%q,", "file_sha256", module.File_sha256)
	json += fmt.Sprintf("%q:%d,", "downloads", module.Downloads)
	json += fmt.Sprintf("%q:%q,", "readme", module.Readme)
	json += fmt.Sprintf("%q:%q,", "changelog", module.Changelog)
	json += fmt.Sprintf("%q:%q,", "license", module.License)
	json += fmt.Sprintf("%q:%q,", "reference", module.Reference)
	json += fmt.Sprintf("%q:%q,", "docs", module.Docs)
	json += fmt.Sprintf("%q:%q,", "pe_compatibility", module.Pe_compatibility)
	json += fmt.Sprintf("%q:%q,", "tasks", module.Tasks)
	json += fmt.Sprintf("%q:%q", "plans", module.Plans)
	json += "}"
	return json
}
