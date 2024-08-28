package sql2json

import (
	"bytes"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

func RowsToJson(rows *sql.Rows) ([]byte, error) {
	if rows == nil {
		return nil, fmt.Errorf("rows is nil")
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var resData []byte
	res := bytes.NewBuffer(resData)
	res.WriteRune('[')

	values, valuesPtr := createPtrs(len(columns))
	for rows.Next() {
		if err := rows.Scan(valuesPtr...); err != nil {
			return nil, err
		}
		res.Write(SerializeRow(columns, values))
		res.WriteRune(',')
	}
	if err := rows.Err(); err != nil {
		res.Reset()
		return nil, err
	}
	if res.Len() > 1 {
		res.Truncate(res.Len() - 1)
	}
	res.WriteRune(']')
	return res.Bytes(), nil
}

func SerializeRow(columns []string, values []sql.RawBytes) []byte {
	var data []byte
	buff := bytes.NewBuffer(data)
	buff.WriteRune('{')
	for i := range columns {
		buff.WriteString(fmt.Sprintf("\"%v\"", spaceToUnderscore(columns[i])))
		buff.WriteRune(':')
		buff.WriteString(parseVal(values[i]))
		if i < len(columns)-1 {
			buff.WriteRune(',')
		}
	}
	buff.WriteRune('}')
	return buff.Bytes()
}

func parseVal(v []byte) string {
	if v == nil {
		return "null"
	}
	// try parse int64
	i64, err := strconv.ParseInt(string(v), 10, 64)
	if err == nil {
		return strconv.FormatInt(i64, 10)
	}
	//try parse float64
	f64, err := strconv.ParseFloat(string(v), 64)
	if err == nil {
		return strconv.FormatFloat(f64, 'f', -1, 64)
	}

	return `"` + string(v) + `"`
}

// Returns slice of pointers
func createPtrs(num int) ([]sql.RawBytes, []interface{}) {
	vals := make([]sql.RawBytes, num)
	ptrs := make([]interface{}, num)
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	return vals, ptrs
}

func spaceToUnderscore(input string) string {
	// Compile the regex to match spaces
	re, err := regexp.Compile(`\s`)
	if err != nil {
		return input
	}
	// Replace all spaces with underscores
	result := re.ReplaceAllString(input, "_")
	return result
}
