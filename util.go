package main

import "reflect"

func contains(key string, slice []string) bool {
	for _, v := range slice {
		if v == key {
			return true
		}
	}
	return false
}

func extractKeys(data []map[string]interface{}) (res []string) {
	for _, m := range data {
		for k, _ := range m {
			if !contains(k, res) {
				res = append(res, k)
			}
		}
	}
	return res
}

func quoteForDriver(driver string) string {
	if driver == "mysql" {
		return "`"
	}
	return `"`
}

func surroundKeysWithQuotes(keys []string, driver string) []string {
	res := make([]string, len(keys))
	quote := quoteForDriver(driver)
	for i, key := range keys {
		res[i] = quote + key + quote
	}
	return res
}

func surroundKeyWithQuote(key, driver string) string {
	quote := quoteForDriver(driver)
	return quote + key + quote
}

// from stretchr/testify/assert/assertions
func ObjectsAreEqual(expected, actual interface{}) bool {

	if expected == nil || actual == nil {
		return expected == actual
	}

	if reflect.DeepEqual(expected, actual) {
		return true
	}

	return false

}

func ObjectsAreEqualValues(expected, actual interface{}) bool {
	if ObjectsAreEqual(expected, actual) {
		return true
	}

	actualType := reflect.TypeOf(actual)
	expectedValue := reflect.ValueOf(expected)
	if expectedValue.Type().ConvertibleTo(actualType) {
		// Attempt comparison after type conversion
		if reflect.DeepEqual(actual, expectedValue.Convert(actualType).Interface()) {
			return true
		}
	}

	return false
}
