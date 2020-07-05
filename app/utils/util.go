package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// NOW generate current datetime
var NOW = time.Now().UTC()

// ReplacePackages --
func ReplacePackages(input string) string {
	paths := strings.Split(input, "/")
	return paths[len(paths)-1]
}

// ParseError --
func ParseError(input string) string {
	fields := []string{}
	errors := strings.Split(input, "\n")
	if len(errors) > 0 {
		for _, err := range errors {
			str := strings.Split(err, "Error:")[0]
			str = strings.ReplaceAll(str, "Key:", "")
			str = strings.ReplaceAll(str, "'", "")
			str = strings.TrimSpace(str)
			str = strings.Split(str, ".")[1]
			fields = append(fields, fmt.Sprintf("`%s`", toSnakeCase(str)))
		}

		return fmt.Sprintf("Field %s tidak boleh kosong.", strings.Join(fields, ", "))
	}

	return input
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// StructToMap get map from struct
func StructToMap(st interface{}) map[string]interface{} {

	reqRules := make(map[string]interface{})

	v := reflect.ValueOf(st)
	t := reflect.TypeOf(st)

	for i := 0; i < v.NumField(); i++ {
		key := strings.ToLower(t.Field(i).Name)
		typ := v.FieldByName(t.Field(i).Name).Kind().String()
		structTag := t.Field(i).Tag.Get("json")
		jsonName := strings.TrimSpace(strings.Split(structTag, ",")[0])
		value := v.FieldByName(t.Field(i).Name)

		// if jsonName is not empty use it for the key
		if jsonName != "" && jsonName != "-" {
			key = jsonName
		}

		if typ == "string" {
			if !(value.String() == "" && strings.Contains(structTag, "omitempty")) {
				fmt.Println(key, value)
				fmt.Println(key, value.String())
				reqRules[key] = value.String()
			}
		} else if typ == "int" {
			reqRules[key] = value.Int()
		} else {
			reqRules[key] = value.Interface()
		}

	}

	return reqRules
}
