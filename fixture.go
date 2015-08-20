package main

import (
	"fmt"
	"strings"
)

type Fixture struct {
	TableName string
	Data      []map[string]interface{}
	Keys      []string
	Update    bool
}

const (
	defaultKey    = "id"
	defaultUpdate = true
)

func typeError(key string, types ...string) error {
	return fmt.Errorf("key '%s' should have one of the following types: %s",
		key, strings.Join(types, ","))
}

func toStringMap(m map[interface{}]interface{}) (map[string]interface{}, error) {
	r := make(map[string]interface{})
	for k, v := range m {
		switch s := k.(type) {
		case string:
			r[s] = v
		default:
			return r, fmt.Errorf("could not convert map key %v to string", k)
		}
	}
	return r, nil
}

func toStringArray(a []interface{}) (as []string, _ error) {
	for _, v := range a {
		switch s := v.(type) {
		case string:
			as = append(as, s)
		default:
			return as, fmt.Errorf("could not convert key %v to string", v)
		}
	}
	return as, nil
}

func convertRecords(m []interface{}) (res []map[string]interface{}, _ error) {
	for _, data := range m {
		switch tmpData := data.(type) {
		case map[string]interface{}:
			res = append(res, tmpData)
		case map[interface{}]interface{}:
			if transformed, err := toStringMap(tmpData); err == nil {
				res = append(res, transformed)
			} else {
				return nil, err
			}
		default:
			return nil, typeError("data", "[]interface{}")
		}
	}
	return res, nil
}

func makeFixture(fixture Fixture, rawData map[string]interface{}) (_ Fixture, err error) {
	data, ok := rawData["data"]
	if !ok {
		return fixture, fmt.Errorf("key 'data' is required")
	}

	if rawInterfaces, ok := data.([]interface{}); ok {
		if fixture.Data, err = convertRecords(rawInterfaces); err != nil {
			return fixture, err
		}
	}

	if rawKeys, ok := rawData["keys"]; ok {
		switch keys := rawKeys.(type) {
		case string:
			fixture.Keys = []string{keys}
		case []string:
			fixture.Keys = keys
		case []interface{}:
			if fixture.Keys, err = toStringArray(keys); err != nil {
				return fixture, err
			}
		default:
			return fixture, typeError("keys", "string", "[]string")
		}
	}

	if rawUpdate, ok := rawData["update"]; ok {
		if update, ok := rawUpdate.(bool); ok {
			fixture.Update = update
		} else {
			return fixture, typeError("update", "bool")
		}
	}

	return fixture, nil
}

func MakeFixture(tableName string, rawData interface{}) (fixture Fixture, err error) {
	fixture = Fixture{
		TableName: tableName,
		Keys:      []string{defaultKey},
		Update:    defaultUpdate,
	}

	switch data := rawData.(type) {
	case []interface{}:
		if fixture.Data, err = convertRecords(data); err != nil {
			return fixture, err
		}
	case map[string]interface{}:
		fixture, err = makeFixture(fixture, data)
	case map[interface{}]interface{}:
		if parsedData, e := toStringMap(data); e != nil {
			err = e
		} else {
			fixture, err = makeFixture(fixture, parsedData)
		}

	default:
		err = fmt.Errorf("can only use array or map as data")
	}

	return fixture, err
}

func MakeFixtures(input map[string]interface{}) (fixtures []Fixture, err error) {
	for tableName, data := range input {
		var fixture Fixture
		if fixture, err = MakeFixture(tableName, data); err != nil {
			return nil, err
		}
		fixtures = append(fixtures, fixture)
	}
	return fixtures, nil
}
