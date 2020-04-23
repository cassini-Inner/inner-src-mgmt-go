package model

import (
	"database/sql"
	"strconv"
)

func ToNullString(value interface{}) sql.NullString{
	if value == nil {
		return sql.NullString{Valid:false}
	}
	switch value.(type) {
	case string: return sql.NullString{
		String: value.(string),
		Valid:  value != "",
	}
	 default:
		return sql.NullString{
			String: strconv.Itoa(int(value.(float64))),
			Valid:  value != "",
		}
	}
}
