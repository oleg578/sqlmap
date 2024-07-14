package sql2json

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
)

func RowsToJson(rows *sql.Rows) ([]byte, error) {
	if rows == nil {
		return nil, fmt.Errorf("rows is nil")
	}
	columns, err := rows.Columns()
	var errSE error
	for i := 0; i < len(columns); i++ {
		columns[i], errSE = spaceToUnderscore(columns[i])
		if errSE != nil {
			return nil, fmt.Errorf("names of columns is invalid: %s", errSE.Error())
		}
	}
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0)
	values, valuePtrs := createPtrs(len(columns))
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		structRow, errSt := slcToStruct(columns, values)
		if errSt != nil {
			panic(errSt)
		}
		result = append(result, structRow.Interface())
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

// Returns slice of pointers
func createPtrs(num int) ([]interface{}, []interface{}) {
	vals := make([]interface{}, num)
	ptrs := make([]interface{}, num)
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	return vals, ptrs
}

// slcToStruct takes in two slices, names and values, and creates a new struct type
// with the specified field names and corresponding types. It returns the created struct type
// as a 'reflect.Type' and an error if the number of names and values are not equal.
func slcToStruct(names []string, values []interface{}) (reflect.Value, error) {
	if len(names) != len(values) {
		return reflect.Value{}, fmt.Errorf("number of names and values are not equal")
	}
	var structFields []reflect.StructField
	for i, name := range names {
		structFields = append(structFields, reflect.StructField{
			Name: name,
			Type: reflect.TypeOf(values[i]),
			Tag:  reflect.StructTag(`json:"` + name + `"`),
		})
	}
	sType := reflect.StructOf(structFields)
	// Create a new instance of the struct
	sValue := reflect.New(sType).Elem()

	// Set the values of the struct fields
	for i, fName := range names {
		sValue.FieldByName(fName).Set(reflect.ValueOf(values[i]))
	}

	return sValue, nil
}

func spaceToUnderscore(input string) (string, error) {
	// Compile the regex to match spaces
	re, err := regexp.Compile(`\s`)
	if err != nil {
		return "", err
	}
	// Replace all spaces with underscores
	result := re.ReplaceAllString(input, "_")
	return result, nil
}
