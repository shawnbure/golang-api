package services

import (
	"errors"
	"fmt"
	"strings"
)

// ConvertFilterToQuery converts a querystring conversion filter to a sql where clause
func ConvertFilterToQuery(tableName string, filter string) (string, []interface{}, error) { //Field|Value|Operator;AND/OR;...
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered in %v", r)
		}
	}()
	var stm string
	var values []interface{}
	if filter == "" {
		return stm, values, nil
	}
	clauses := strings.Split(filter, ";")
	for _, c := range clauses {
		if c == "OR" || c == "AND" {
			stm += " " + c + " "
			continue
		}
		params := strings.Split(c, "|")
		// We must have only 3 length of params. field|value|operator
		if len(params) != 3 {
			fmt.Printf("we have wrong params structure '%s' in '%s'", c, filter)
			continue
		}
		field := params[0]
		value := params[1]
		operator := params[2]
		prefix := tableName
		subObjects := strings.Split(field, ".")

		if len(subObjects) > 1 {
			prefix = "\"" + subObjects[0] + "\""
			field = strings.Join(subObjects[1:], ".")
		}

		var query string

		if operator == "BETWEEN" {
			ranges := strings.Split(value, "AND")
			if len(ranges) != 2 {
				err := errors.New("bad given between range")
				return "", nil, err
			}
			query = prefix + "." + field + " BETWEEN " + "?" + " AND " + "?"
			values = append(values, ranges[0])
			values = append(values, ranges[1])
		} else {
			query = prefix + "." + field + " " + operator + " " + "?"
			values = append(values, value)
		}
		stm += query
	}
	return stm, values, nil
}
