package notify

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten/internal/common"
)

func (l *Listener) EnsureNotifyInfra(ctx context.Context) error {
	_, err := l.pool.Exec(ctx, `
CREATE OR REPLACE FUNCTION public.notify_row_change() RETURNS trigger AS $$
DECLARE payload jsonb;
BEGIN
  payload := jsonb_build_object(
    'table', TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME,
    'op', TG_OP,
    'new', CASE WHEN TG_OP IN ('INSERT','UPDATE') THEN to_jsonb(NEW) ELSE NULL END,
    'old', CASE WHEN TG_OP IN ('UPDATE','DELETE') THEN to_jsonb(OLD) ELSE NULL END
  );
  PERFORM pg_notify('row_changes', payload::text);
  RETURN COALESCE(NEW, OLD);
END; $$ LANGUAGE plpgsql;`)
	return err
}

func (l *Listener) EnsureNotifyTrigger(ctx context.Context, schema, table string) error {
	trg := fmt.Sprintf("trg_%s_row_notify", table) // simple derived name

	sql := fmt.Sprintf(`
DO $$
DECLARE trg text := %s;
DECLARE sch text := %s;
DECLARE tbl text := %s;
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_trigger t
    JOIN pg_class c ON c.oid = t.tgrelid
    JOIN pg_namespace n ON n.oid = c.relnamespace
    WHERE t.tgname = trg AND n.nspname = sch AND c.relname = tbl
  ) THEN
    EXECUTE format('CREATE TRIGGER %%I
      AFTER INSERT OR UPDATE OR DELETE ON %%I.%%I
      FOR EACH ROW EXECUTE PROCEDURE public.notify_row_change()', trg, sch, tbl);
  END IF;
END $$;`,
		common.QuoteLit(trg), common.QuoteLit(schema), common.QuoteLit(table),
	)

	_, err := l.pool.Exec(ctx, sql)
	return err
}
