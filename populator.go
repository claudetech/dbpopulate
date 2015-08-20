package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type Populator struct {
	DB                      *sql.DB
	Driver                  string
	currentPlaceholderIndex int
}

func (p *Populator) getPlaceholder() string {
	if p.Driver == "postgres" {
		val := p.currentPlaceholderIndex
		p.currentPlaceholderIndex++
		return fmt.Sprintf("$%d", val)
	}
	return "?"
}

func (p *Populator) resetPlaceholder() {
	p.currentPlaceholderIndex = 1
}

func (p *Populator) PopulateData(fixtures []Fixture) error {
	for _, f := range fixtures {
		if err := p.PopulateFixture(f); err != nil {
			return err
		}
	}
	return nil
}

func (p *Populator) generateCondition(keys []string, data map[string]interface{}) (string, []interface{}, error) {
	var conditions []string
	var args []interface{}
	for _, key := range keys {
		if v, ok := data[key]; ok {
			conditions = append(conditions, fmt.Sprintf("%s=%s", key, p.getPlaceholder()))
			args = append(args, v)
		} else {
			return "", nil, fmt.Errorf("key %s not found in record %v", key, data)
		}
	}
	return "(" + strings.Join(conditions, " AND ") + ")", args, nil
}

func (p *Populator) generateSelectQuery(fixture Fixture) (string, []interface{}, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE ", strings.Join(fixture.Keys, ", "), fixture.TableName)
	var conditions []string
	var args []interface{}
	for _, data := range fixture.Data {
		condition, cArgs, err := p.generateCondition(fixture.Keys, data)
		if err != nil {
			return "", nil, err
		}
		conditions = append(conditions, condition)
		args = append(args, cArgs...)
	}
	return query + strings.Join(conditions, " OR "), args, nil
}

func (p *Populator) makePlaceholders(args []interface{}) string {
	var placeholders []string
	for _, _ = range args {
		placeholders = append(placeholders, p.getPlaceholder())
	}
	return "(" + strings.Join(placeholders, ",") + ")"
}

func (p *Populator) generateInsertStmt(fixture Fixture, data []map[string]interface{}) (string, []interface{}) {
	query := "INSERT INTO %s (%s) VALUES %s"
	keys := extractKeys(data)
	var args []interface{}
	var placeholders []string
	for _, record := range data {
		var recordArgs []interface{}
		for _, k := range keys {
			if v, ok := record[k]; ok {
				recordArgs = append(recordArgs, v)
			}
		}
		placeholders = append(placeholders, p.makePlaceholders(recordArgs))
		args = append(args, recordArgs...)
	}
	values := strings.Join(placeholders, ",")
	return fmt.Sprintf(query, fixture.TableName, strings.Join(keys, ","), values), args
}

func (p *Populator) getExistingData(fixture Fixture) ([]map[string]interface{}, error) {
	query, args, err := p.generateSelectQuery(fixture)
	if err != nil {
		return nil, err
	}

	logger.Debugf("executing query %s with args %+v", query, args)
	rows, err := p.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	p.resetPlaceholder()

	var result []map[string]interface{}

	values := make([]interface{}, len(fixture.Keys))
	valuesPtr := make([]interface{}, len(fixture.Keys))

	for rows.Next() {
		row := make(map[string]interface{})
		for i := 0; i < len(fixture.Keys); i++ {
			valuesPtr[i] = &values[i]
		}
		rows.Scan(valuesPtr...)
		for i, res := range values {
			row[fixture.Keys[i]] = res
		}
		result = append(result, row)
	}

	return result, nil
}

func (p *Populator) hasRecord(keys []string, existingData []map[string]interface{}, record map[string]interface{}) bool {
CheckEquality:
	for _, data := range existingData {
		for _, k := range keys {
			if !ObjectsAreEqualValues(data[k], record[k]) {
				continue CheckEquality
			}
		}
		return true
	}
	return false
}

func (p *Populator) getNewData(fixture Fixture, existingData []map[string]interface{}) (newData []map[string]interface{}) {
	for _, record := range fixture.Data {
		if !p.hasRecord(fixture.Keys, existingData, record) {
			newData = append(newData, record)
		}
	}
	return newData
}

func (p *Populator) PopulateFixture(fixture Fixture) error {
	existingData, err := p.getExistingData(fixture)
	if err != nil {
		return err
	}
	newData := p.getNewData(fixture, existingData)
	if len(newData) == 0 {
		return nil
	}
	insertStmt, args := p.generateInsertStmt(fixture, newData)
	logger.Infof("inserting %d new records in %s", len(newData), fixture.TableName)
	logger.Debugf("executing %s with args %+v", insertStmt, args)
	_, err = p.DB.Exec(insertStmt, args...)
	p.resetPlaceholder()
	return err
}

func NewPopulator(dbUrl string) (*Populator, error) {
	db, driverName, err := ConnectToDb(dbUrl)
	if err != nil {
		return nil, err
	}
	return &Populator{
		DB:                      db,
		Driver:                  driverName,
		currentPlaceholderIndex: 1,
	}, nil
}
