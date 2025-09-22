package queue

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten/internal/common"
)

func (l *Listener) EnsureQueueInfra(ctx context.Context) error {
	_, err := l.pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS public.row_change_queue (
  id           bigserial PRIMARY KEY,
  tbl          text      NOT NULL,
  op           text      NOT NULL,
  new_json     jsonb,
  old_json     jsonb,
  created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE OR REPLACE FUNCTION public.enqueue_row_change_pk() RETURNS trigger AS $$
BEGIN
  INSERT INTO public.row_change_queue(row_id, tbl, op, new_json, old_json)
  VALUES (
    CASE 
      WHEN TG_OP IN ('UPDATE','DELETE') THEN OLD.id
      ELSE NEW.id
    END,
    TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME,
    TG_OP,
    CASE WHEN TG_OP IN ('INSERT','UPDATE') THEN to_jsonb(NEW) ELSE NULL END,
    CASE WHEN TG_OP IN ('UPDATE','DELETE') THEN to_jsonb(OLD) ELSE NULL END
  );
  PERFORM pg_notify('row_change_queue', '');
  RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;`)
	return err
}

func (l *Listener) EnsureQueueTrigger(ctx context.Context, schema, table string) error {
	trg := fmt.Sprintf("trg_%s_row_queue", table)

	sql := fmt.Sprintf(`
DO $do$
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
      FOR EACH ROW EXECUTE PROCEDURE public.enqueue_row_change_pk()', trg, sch, tbl);
  END IF;
END
$do$;`,
		common.QuoteLit(trg), common.QuoteLit(schema), common.QuoteLit(table),
	)

	_, err := l.pool.Exec(ctx, sql)
	return err
}
