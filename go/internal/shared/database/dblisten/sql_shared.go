package dblisten

import (
	"context"
	"fmt"
	"strings"
)

// sql_shared.go
func (l *Listener) ensureTrigger(ctx context.Context, schema, table, proc string) error {
	// proc must look like: "public.notify_row_change()"
	trg := fmt.Sprintf("trg_%s_%s",
		strings.ReplaceAll(table, `"`, ``),
		strings.TrimSuffix(strings.Split(strings.TrimPrefix(proc, "public."), ".")[len(strings.Split(strings.TrimPrefix(proc, "public."), "."))-1], "()"),
	)

	sql := fmt.Sprintf(`
DO $do$
DECLARE
  trg text := %s;
  sch text := %s;
  tbl text := %s;
  prc text := %s;
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_trigger t
    JOIN pg_class c ON c.oid = t.tgrelid
    JOIN pg_namespace n ON n.oid = c.relnamespace
    WHERE t.tgname = trg AND n.nspname = sch AND c.relname = tbl
  ) THEN
    EXECUTE
      'CREATE TRIGGER ' || quote_ident(trg) ||
      ' AFTER INSERT OR UPDATE OR DELETE ON ' || quote_ident(sch) || '.' || quote_ident(tbl) ||
      ' FOR EACH ROW EXECUTE PROCEDURE ' || prc;
  END IF;
END
$do$;`,
		quoteLit(trg), quoteLit(schema), quoteLit(table), quoteLit(proc),
	)

	_, err := l.conn.Exec(ctx, sql)
	return err
}
