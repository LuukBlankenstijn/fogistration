-- Write your migrate up statements here

CREATE OR REPLACE FUNCTION notify_ip_change() RETURNS trigger AS $$
DECLARE
    data JSONB;
    notify_type TEXT := 'change_ip';
BEGIN
    IF NEW.ip IS DISTINCT FROM OLD.ip THEN
        data := jsonb_build_object(
            'id', NEW.id,
            'ip', to_jsonb(NEW.ip),
            'ip_old', to_jsonb(OLD.ip)
        );

        -- Send NOTIFY with the same type
        PERFORM pg_notify(notify_type, data::text);

        -- Insert into command_queue with the same type
        INSERT INTO message_queue (command_type, payload)
        VALUES (notify_type, data);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_ip_change
AFTER UPDATE ON teams
FOR EACH ROW
WHEN (OLD.ip IS DISTINCT FROM NEW.ip)
EXECUTE FUNCTION notify_ip_change();

---- create above / drop below ----

DROP TRIGGER IF EXISTS trigger_ip_change on teams;
DROP FUNCTION IF EXISTS notify_ip_change();

