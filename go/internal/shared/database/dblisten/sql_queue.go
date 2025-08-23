package dblisten

import (
	"context"
	"fmt"
	"strings"
)

func (l *Listener) EnsureQueueInfra(ctx context.Context) error {
	_, err := l.conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS public.row_change_queue (
  id           bigserial PRIMARY KEY,
  tbl          text        NOT NULL,
  row_id       text        NOT NULL,
  op           text        NOT NULL,
  new_json     jsonb,
  old_json     jsonb,
  enqueued_at  timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  UNIQUE (tbl, row_id)
);
CREATE INDEX IF NOT EXISTS row_change_queue_tbl_rowid_idx ON public.row_change_queue(tbl, row_id);

CREATE OR REPLACE FUNCTION public.enqueue_row_change_pk() RETURNS trigger AS $$
DECLARE
  v_new    jsonb;
  v_old    jsonb;
  v_row_id text;
  v_table  text := TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
BEGIN
  -- fixed PK column 'id'
  IF TG_OP = 'DELETE' THEN
    EXECUTE 'SELECT ($1).id::text' INTO v_row_id USING OLD;
  ELSE
    EXECUTE 'SELECT ($1).id::text' INTO v_row_id USING NEW;
  END IF;

  v_new := CASE WHEN TG_OP IN ('INSERT','UPDATE') THEN to_jsonb(NEW) ELSE NULL END;
  v_old := CASE WHEN TG_OP IN ('UPDATE','DELETE') THEN to_jsonb(OLD) ELSE NULL END;

  INSERT INTO public.row_change_queue (tbl, row_id, op, new_json, old_json)
  VALUES (v_table, v_row_id, TG_OP, v_new, v_old)
  ON CONFLICT (tbl, row_id) DO UPDATE SET
    op         = EXCLUDED.op,
    new_json   = COALESCE(EXCLUDED.new_json, public.row_change_queue.new_json),
    old_json   = COALESCE(EXCLUDED.old_json, public.row_change_queue.old_json),
    updated_at = now();

  PERFORM pg_notify('row_change_queue', v_table || ':' || v_row_id);
  RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;`)
	return err
}

func (l *Listener) EnsureQueueTrigger(ctx context.Context, schema, table string) error {
	trg := fmt.Sprintf("trg_%s_row_queue", strings.ReplaceAll(table, `"`, ``))

	sql := fmt.Sprintf(`
DO $do$
DECLARE
  trg text := %s;
  sch text := %s;
  tbl text := %s;
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
      ' FOR EACH ROW EXECUTE PROCEDURE public.enqueue_row_change_pk()';
  END IF;
END
$do$;`,
		quoteLit(trg), quoteLit(schema), quoteLit(table),
	)

	_, err := l.conn.Exec(ctx, sql)
	return err
}
