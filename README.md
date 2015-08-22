# dbpopulate [![Build Status](https://travis-ci.org/claudetech/dbpopulate.svg?branch=master)](https://travis-ci.org/claudetech/dbpopulate)

CLI tool written in Go to populate an SQL database from JSON or YAML data.

## Installation

```sh
$ go get github.com/claudetech/dbpopulate
```

## Usage

The only required option is `--database-url`.
It can be of the form `postgres://POSTGRES_URL`, `mysql://MYSQL_URL` or `sqlite3://PATH_TO_DB`.
You can also use the `DATABASE_URL` environment variable instead.

```sh
$ dbpopulate --debug --env=development --database-url=sqlite3://mydata.db
$ dbpopulate --quiet --database-url=postgres://localhost/foobar?sslmode=disable
$ dbpopulate --database-url=mysql://foobar:password@tcp(localhost:3306)/foobar --fixtures-path=/path/to/my/fixtures
```

## Fixture files

By default, `dbpopulate` will read all the files ending in `.yml`, `.yaml` or `.json` present in `FIXTURES_PATH` and `FIXTURES_PATH/$GO_ENV` directory. If the latter does not exist, it will be ignored.

`FIXTURES_PATH` defaults to `./fixtures` and can be changed using the `--fixtures-path` option.
You can change `GO_ENV` by setting the `GO_ENV` environment variable or with the `--env` flag.

Here is a sample fixture file.

```yaml
countries:
  - id: 1
    name: 'France'
  - id: 2
    name: 'Japan'

users:
  - id: 1
    name: 'tuvistavie'
    country_id: 1
```

where each key is a table name, and each value are the records to add.
If you want to avoid passing the `id` and use another unique key for the records, you can use the following form:

```yaml
countries:
  keys: [name]
  data:
      - name: 'France'
      - name: 'Japan'
```

You can use a single, or multiple keys to distinguish the records.

Here is an example in JSON:

```json
{
  "regions": {
    "keys": ["name", "country_id"],
    "data": [{
      "country_id": 1,
      "name": "Ile de france"
    }, {
      "country_id": 2,
      "name": "Tokyo"
    }]
  },
  "prefectures": [{
    "id": 1,
    "region_id": 2,
    "name": "千代田"
  }]
}
```

## CLI options

The available CLI options (and there environment variable equivalent) are:

* `--database-url` (`-u`, `$DATABASE_URL`): Database URL
* `--fixtures-path` (`-p`, `$FIXTURES_PATH`): Path to the directory containing fixtures
* `--env` (`-e`, `$GO_ENV`): Environment (used to look for subdirectories)
* `--debug` (`-d`, `$DEBUG`): Activate debug mode (more log)
* `--quiet` (`-q`, `$QUIET`): Activate quiet mode (less log)

dbpopulate uses [dotenv](https://github.com/joho/godotenv) to load environment variables, so you can put a `.env` file at the top of your project with the needed settings and use the `dbpopulate` command without any options.

## Contributing

Please feel free to add support for other DB drivers,
or other seed files format if you need.

To increasing logging level, you can pass the `--debug` flag or set the
`DEBUG` environment to anything not empty.
