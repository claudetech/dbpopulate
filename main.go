package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/claudetech/loggo"
	"github.com/codegangsta/cli"
)

func showErrorAndAbort(mesasge string) {
	fmt.Fprintf(os.Stderr, "An error has occured: %s, aborting.\n", mesasge)
	os.Exit(1)
}

func GetFixturesDirectories(fixturesPath string, env string) []string {
	fixturesPaths := []string{fixturesPath}
	if env != "" {
		envDir := filepath.Join(fixturesPath, env)
		if _, err := os.Stat(envDir); err == nil {
			fixturesPaths = append(fixturesPaths, envDir)
		}
	}
	return fixturesPaths
}

func RunApp(c *cli.Context) {
	dbUrl := c.String("db-url")

	if dbUrl == "" {
		showErrorAndAbort("you need to provide 'db-url' or $DATABASE_URL needs to be set")
	}

	if c.Bool("debug") {
		logger.SetLevel(loggo.Debug)
	}

	if c.Bool("quiet") {
		logger.SetLevel(loggo.Warning)
	}

	fixturesPaths := GetFixturesDirectories(c.String("fixtures-path"), c.String("env"))
	logger.Debugf("searching for fixtures in: %v", fixturesPaths)

	fixtures, err := LoadDirectories(fixturesPaths)
	if err != nil {
		showErrorAndAbort(err.Error())
	}
	logger.Debugf("found the following fixtures: %+v", fixtures)

	logger.Infof("connecting to DB %s", dbUrl)

	populator, err := NewPopulator(dbUrl)
	if err != nil {
		showErrorAndAbort(err.Error())
	}

	if err := populator.PopulateData(fixtures); err != nil {
		showErrorAndAbort(err.Error())
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "dbpopulate"
	app.Usage = "populates SQL database from JSON or YAML files"
	app.Author = "Daniel Perez <daniel@claudetech.com>"
	app.Version = "0.1.0"
	app.Action = RunApp
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "fixtures-path, p",
			Value: "./fixtures",
			Usage: "Set the directory containing the fixtures",
		},
		cli.StringFlag{
			Name:   "db-url, u",
			Usage:  "Set the database URL",
			EnvVar: "DATABASE_URL",
		},
		cli.StringFlag{
			Name:   "env, e",
			Usage:  "Set the environment",
			EnvVar: "GO_ENV",
		},
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "Set debug mode on",
			EnvVar: "DEBUG",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Disable info log",
		},
	}
	app.Run(os.Args)
}
