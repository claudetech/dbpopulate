package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	gzipExt = ".gz"
)

var availableExts = [...]string{".json", ".yml", ".yaml"}

var mappings = map[string]string{
	".yaml": ".yml",
}

type FileLoader func(string) (map[string]interface{}, error)

func isAvailableExt(ext string) bool {
	for _, e := range availableExts {
		if e == ext {
			return true
		}
	}
	return false
}

func uncompressData(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func GetNormalizedExtension(path string) string {
	ext := filepath.Ext(path)
	if ext == gzipExt {
		ext = filepath.Ext(strings.TrimSuffix(path, ext))
	}
	if mapping, ok := mappings[ext]; ok {
		ext = mapping
	}
	return ext
}

func MakeLoader(loader func([]byte, interface{}) error) FileLoader {
	return func(path string) (output map[string]interface{}, _ error) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return output, err
		}
		if filepath.Ext(path) == gzipExt {
			if data, err = uncompressData(data); err != nil {
				return output, err
			}
		}
		err = loader(data, &output)
		return output, err
	}
}

func MakeLoaderFor(path string) (FileLoader, error) {
	ext := GetNormalizedExtension(path)
	switch ext {
	case ".json":
		return MakeLoader(json.Unmarshal), nil
	case ".yml":
		return MakeLoader(yaml.Unmarshal), nil
	default:
		return nil, fmt.Errorf("no loader available for file with extension %s", ext)
	}
}

func LoadFile(path string) ([]Fixture, error) {
	load, err := MakeLoaderFor(path)
	if err != nil {
		return nil, err
	}
	content, err := load(path)
	if err != nil {
		return nil, err
	}
	return MakeFixtures(content)
}

func getFileNames(files []os.FileInfo) (names []string) {
	for _, file := range files {
		names = append(names, file.Name())
	}
	sort.Strings(names)
	return names
}

func LoadDirectory(directory string) (fixtures []Fixture, err error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("could not read directory %s", directory)
	}
	fileNames := getFileNames(files)
	for _, fileName := range fileNames {
		if isAvailableExt(GetNormalizedExtension(fileName)) {
			content, err := LoadFile(filepath.Join(directory, fileName))
			if err != nil {
				return nil, err
			}
			fixtures = append(fixtures, content...)
		}
	}

	return fixtures, nil
}

func LoadDirectories(directories []string) (fixtures []Fixture, err error) {
	for _, directory := range directories {
		results, err := LoadDirectory(directory)
		if err != nil {
			return nil, err
		}
		fixtures = append(fixtures, results...)
	}
	return fixtures, nil
}
