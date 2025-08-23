package dblisten

import "context"

func (l *Listener) EnsureNotifyInfra(ctx context.Context) error {
	_, err := l.conn.Exec(ctx, `
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
END;
$$ LANGUAGE plpgsql;`)
	return err
}

func (l *Listener) EnsureNotifyTrigger(ctx context.Context, schema, table string) error {
	return l.ensureTrigger(ctx, schema, table, "public.notify_row_change()")
}
