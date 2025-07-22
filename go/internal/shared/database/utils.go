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

func PgTextFromString(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}
