package sql2json

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func RowsToJson(rows *sql.Rows) ([]byte, error) {
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
			colName := normalizeName(col)
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

// normalizeName Replaces spaces with underscores in the column names and capitalizes the first letter.
func normalizeName(input string) string {
	re := regexp.MustCompile(`\s`)
	// capitalize first symbol of input
	input = strings.ToUpper(input[:1]) + input[1:]
	return re.ReplaceAllString(input, "_")
}
