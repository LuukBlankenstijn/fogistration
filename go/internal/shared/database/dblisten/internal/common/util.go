package common

import "strings"

func NormalizeTable(table string) (string, string) {
	s := "public"
	t := table
	if strings.Contains(table, ".") {
		parts := strings.SplitN(table, ".", 2)
		s, t = parts[0], parts[1]
	}
	return s, t
}

func ToKey(table string) string {
	s, t := NormalizeTable(table)
	return s + "." + t
}

func QuoteLit(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }
func QuoteIdent(s string) string { return `"` + strings.ReplaceAll(s, `"`, `""`) + `"` }
