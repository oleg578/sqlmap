package sql2json

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
)

func RowsToJsonMapping(rows *sql.Rows) ([]byte, error) {
	if rows == nil {
		return nil, fmt.Errorf("rows is nil")
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0)
	values, valuePtrs := createPtrs(len(columns))
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			colName := spaceToUnderscore(col)
			rowMap[colName] = assignCellValue(values[i])
		}
		result = append(result, rowMap)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

func assignCellValue(val interface{}) interface{} {
	if b, ok := val.([]byte); ok {
		return string(b)
	}
	return val
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

// RowsToJson converts the result of a SQL query (sql.Rows) into a JSON encoded byte array.
func RowsToJsonReflection(rows *sql.Rows) ([]byte, error) {
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

// spaceToUnderscore replaces all space characters in the input string with underscores.
func spaceToUnderscore(input string) string {
	re := regexp.MustCompile(`\s`)
	return re.ReplaceAllString(input, "_")
}

// buildStruct takes a slice of SQL column types and
// returns a new 'reflect.Type' representing a struct with fields
// corresponding to the columns, where the field names are converted from spaces
// to underscores and tagged with json tags.
// Example: buildStruct(columns []*sql.ColumnType) reflect.Type
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
