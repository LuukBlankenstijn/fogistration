package database

import (
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func StringValueFromPgText(text pgtype.Text) *wrapperspb.StringValue {
	if !text.Valid {
		return nil
	}
	return wrapperspb.String(text.String)
}

func PgTextFromString(value *string) pgtype.Text {
	var v string
	if value != nil {
		v = *value
	} else {
		v = ""
	}
	return pgtype.Text{
		String: v,
		Valid:  value != nil,
	}
}
