package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	dbURL              = "sqlite3://:memory:"
	schemaFileTemplate = "./db/tables.%s.sql"
)

var populator *Populator

func assertCount(t *testing.T, table string, length int) {
	res, err := populator.DB.Query("SELECT COUNT(*) FROM " + table)
	assert.Nil(t, err)
	defer res.Close()
	for res.Next() {
		var count int
		assert.Nil(t, res.Scan(&count))
		assert.Equal(t, length, count)
	}
}

func showErrAndExitTest(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func TestMain(m *testing.M) {
	var err error

	connectionURL := dbURL
	if envURL := os.Getenv("DATABASE_URL"); envURL != "" {
		connectionURL = envURL
	}

	if populator, err = NewPopulator(connectionURL); err != nil {
		showErrAndExitTest(err)
	}

	uri, err := url.Parse(connectionURL)
	if err != nil {
		showErrAndExitTest(err)
	}

	schemaFile := fmt.Sprintf(schemaFileTemplate, uri.Scheme)
	createTablesStmt, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		showErrAndExitTest(err)
	}

	statements := strings.Split(string(createTablesStmt), ";")
	for _, statement := range statements {
		if strings.TrimSpace(statement) == "" {
			continue
		}
		if _, err = populator.DB.Exec(statement); err != nil {
			showErrAndExitTest(err)
		}
	}

	os.Exit(m.Run())
}

func TestPopulateFixture(t *testing.T) {
	fixtures, err := LoadFile("./fixtures/001_countries.yml")
	assert.Nil(t, err)
	err = populator.PopulateFixture(fixtures[0])
	assert.Nil(t, err)
	assertCount(t, "countries", 2)
}

func TestPopulateData(t *testing.T) {
	fixtures, err := LoadDirectory("./fixtures")
	assert.Nil(t, err)
	err = populator.PopulateData(fixtures)
	assert.Nil(t, err)
	assertCount(t, "countries", 2)
	assertCount(t, "regions", 2)
	assertCount(t, "prefectures", 1)
}
