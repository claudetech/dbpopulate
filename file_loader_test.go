package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNormalizedExtension(t *testing.T) {
	for input, expected := range map[string]string{
		"foo.json":    ".json",
		"foo.yml":     ".yml",
		"foo.yaml":    ".yml",
		"foo.yaml.gz": ".yml",
		"foo.json.gz": ".json",
	} {
		assert.Equal(t, expected, GetNormalizedExtension(input))
	}
}

func TestMakeLoaderFor(t *testing.T) {
	for input, shouldHaveError := range map[string]bool{
		"path/to/myfile.json":   false,
		"path/to/myfile.yml":    false,
		"path/to/myfile.yml.gz": false,
		"path/to/myfile.sql":    true,
	} {
		_, err := MakeLoaderFor(input)
		if shouldHaveError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestLoadFile(t *testing.T) {
	_, err := LoadFile("/path/to/inexisting")
	assert.NotNil(t, err)
	for _, testCase := range []struct {
		FileName    string
		FixturesNum int
		Data        map[string]int
	}{
		{"./fixtures/001_countries.yml", 1, map[string]int{"countries": 2}},
		{"./fixtures/002_regions.json.gz", 2, map[string]int{"prefectures": 1, "regions": 2}},
	} {
		fixtures, err := LoadFile(testCase.FileName)
		assert.Nil(t, err)
		assert.Len(t, fixtures, testCase.FixturesNum)
		for _, fixture := range fixtures {
			assert.Len(t, fixture.Data, testCase.Data[fixture.TableName])
		}
	}
}

func TestLoadDirectory(t *testing.T) {
	_, err := LoadDirectory("/path/to/inexisting")
	assert.NotNil(t, err)
	fixtures, err := LoadDirectory("./fixtures")
	assert.Nil(t, err)
	assert.Len(t, fixtures, 3)
}
