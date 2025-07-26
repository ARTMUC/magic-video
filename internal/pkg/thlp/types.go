package thlp

import (
	"database/sql"
)

func PatchStringToString(in string, out string) string {
	if in == "" {
		return out
	}
	return in
}

func SqlNullstringToPtrString(in sql.Null[string]) *string {
	if in.Valid {
		return &in.V
	}
	return nil
}
