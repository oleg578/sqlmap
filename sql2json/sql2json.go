package sql2json

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type Dummy struct {
	ID       int64
	Product  string
	Price    float64
	Qty      int64
	NullData sql.NullString
	Date     string
}

func RowsToJson(rows *sql.Rows) ([]byte, error) {
	var result = make([]Dummy, 0)
	if rows == nil {
		return nil, errors.New("rows is nil")
	}
	for rows.Next() {
		var rec = Dummy{}
		if err := rows.Scan(
			&rec.ID,
			&rec.Product,
			&rec.Price,
			&rec.Qty,
			&rec.NullData,
			&rec.Date); err != nil {
			return nil, err
		}
		if !rec.NullData.Valid {
			rec.NullData.String = "null"
		}
		result = append(result, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return json.Marshal(result)
}
