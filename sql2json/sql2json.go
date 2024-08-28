package sql2json

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strings"
)

// RowsToJson converts the result of a SQL query (sql.Rows) into a JSON encoded byte array.
func RowsToJson(rows *sql.Rows) (<-chan interface{}, <-chan error) {
	ch := make(chan interface{})
	chErr := make(chan error)
	go func() {
		defer close(ch)
		defer close(chErr)
		if rows == nil {
			chErr <- errors.New("rows is nil")
			return
		}
		columns, err := rows.ColumnTypes()
		if err != nil {
			chErr <- err
			return
		}
		if len(columns) == 0 {
			chErr <- errors.New("no columns found")
		}
		num := len(columns)
		values, valuePtrs := make([]interface{}, num), make([]interface{}, num)
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		structRow := buildStruct(columns)
		for rows.Next() {
			if err := rows.Scan(valuePtrs...); err != nil {
				chErr <- err
				return
			}
			rowVal := reflect.New(structRow).Elem()
			// Set the values of the struct fields
			for i, col := range columns {
				val := reflect.ValueOf(*(valuePtrs[i].(*interface{})))
				fieldName := normalizeName(col.Name())
				field := rowVal.FieldByName(fieldName)
				if field.IsValid() {
					if val.IsValid() {
						switch col.ScanType().Name() {
						case "string":
							field.SetString(string(val.Bytes()))
						default:
							field.Set(reflect.ValueOf(val.Interface()))
						}
					} else {
						field.SetZero()
					}
				}
			}
			ch <- rowVal.Interface()
		}
		if err := rows.Err(); err != nil {
			chErr <- err
		}
	}()
	return ch, chErr
}

// normalizeName Replaces spaces with underscores in the column names and capitalizes the first letter.
func normalizeName(input string) string {
	re := regexp.MustCompile(`\s`)
	// capitalize first symbol of input
	input = strings.ToUpper(input[:1]) + input[1:]
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
		name := normalizeName(col.Name())
		structFields = append(structFields, reflect.StructField{
			Name: name,
			Type: col.ScanType(),
			Tag:  reflect.StructTag(`json:"` + name + `"`),
		})
	}
	return reflect.StructOf(structFields)
}
