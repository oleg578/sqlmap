package sql2json

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
)

func RowsToJson(rows *sql.Rows) ([]byte, error) {
	if rows == nil {
		return nil, fmt.Errorf("rows is nil")
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	values, valuePtrs := createPtrs(len(columns))
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			colName, err := spaceToUnderscore(col)
			if err != nil {
				return nil, err
			}
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
