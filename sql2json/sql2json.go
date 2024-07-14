package sql2json

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func RowsToJson(rows *sql.Rows) ([]byte, error) {
	if rows == nil {
		return nil, fmt.Errorf("rows is nil")
	}
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("no columns found")
	}
	result := make([]interface{}, 0)
	values, valuePtrs := createPtrs(len(columns))
	structRow := buildStruct(columns)
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		rowVal := reflect.New(structRow).Elem()
		// Set the values of the struct fields
		for i, col := range columns {
			rowVal.FieldByName(spaceToUnderscore(col.Name())).Set(reflect.ValueOf(values[i]))
		}
		result = append(result, rowVal.Interface())
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

func spaceToUnderscore(input string) string {
	return strings.Replace(input, " ", "_", -1)
}

// buildStruct takes in a map of field names and their corresponding types, and creates a new struct type
// with the specified field names and types. It returns the created struct type as a 'reflect.Type'.
// The field names are used as the struct field names, and the field types are used as the struct field types.
// The struct field tags are set with `json` tags using the field names.
func buildStruct(columns []*sql.ColumnType) reflect.Type {
	var structFields []reflect.StructField
	for _, col := range columns {
		name := spaceToUnderscore(col.Name())
		structFields = append(structFields, reflect.StructField{
			Name: name,
			Type: col.ScanType(),
			Tag:  reflect.StructTag(`json:"` + name + `"`),
		})
	}
	return reflect.StructOf(structFields)
}
