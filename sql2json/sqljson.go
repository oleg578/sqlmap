package sql2json

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Dummy struct {
	Column1 string
	Column2 string
}

func RowsToJson(rows *sql.Rows) ([]byte, error) {
	var result = make([]Dummy, 0)
	if rows == nil {
		return nil, fmt.Errorf("rows is nil")
	}
	var rec = Dummy{}
	for rows.Next() {
		if err := rows.Scan(&rec.Column1, &rec.Column2); err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return json.Marshal(result)
}
