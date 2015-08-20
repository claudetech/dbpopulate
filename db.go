package main

import (
	"database/sql"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func normalizeDbUrl(provider string, dbUrl string) string {
	switch provider {
	case "sqlite3", "mysql":
		return strings.TrimPrefix(dbUrl, provider+"://")
	default:
		return dbUrl
	}
	return dbUrl
}

func dbDataFromUrl(dbConnectionString string) (string, string, error) {
	dbUrl, err := url.Parse(dbConnectionString)
	if err != nil {
		return "", "", err
	}
	return dbUrl.Scheme, normalizeDbUrl(dbUrl.Scheme, dbUrl.String()), nil
}

func ConnectToDb(dbConnectionString string) (*sql.DB, string, error) {
	driverName, dataSourceName, err := dbDataFromUrl(dbConnectionString)
	if err != nil {
		return nil, "", err
	}
	db, err := sql.Open(driverName, dataSourceName)
	return db, driverName, err
}
