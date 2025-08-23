package dblisten

import "strings"

func normalizeTable(in string) (schema, table string) {
	in = strings.TrimSpace(in)
	if in == "" {
		return "public", ""
	}
	parts := strings.Split(in, ".")
	if len(parts) == 1 {
		return "public", parts[0]
	}
	return parts[0], parts[1]
}

func toKey(table string) string {
	s, t := normalizeTable(table)
	return s + "." + t
}

func quoteIdent(id string) string {
	return `"` + strings.ReplaceAll(id, `"`, `""`) + `"`
}

func quoteLit(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `''`) + `'`
}
