package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	dbFile     = "db/data.sqlite.db"
	dbUrl      = "sqlite3://" + dbFile
	schemaFile = "./db/tables.sqlite3.sql"
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
	os.Remove(dbFile)
	var err error
	if populator, err = NewPopulator(dbUrl); err != nil {
		showErrAndExitTest(err)
	}
	createTablesStmt, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		showErrAndExitTest(err)
	}
	if _, err = populator.DB.Exec(string(createTablesStmt)); err != nil {
		showErrAndExitTest(err)
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
