-- Write your migrate up statements here
CREATE TABLE message_queue (
    id SERIAL PRIMARY KEY,
    command_type VARCHAR NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Trigger for notifications
CREATE OR REPLACE FUNCTION notify_new_command()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('new_command', '');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_notify_new_command
    AFTER INSERT ON message_queue
    FOR EACH ROW
    EXECUTE FUNCTION notify_new_command();

---- create above / drop below ----
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

DROP TRIGGER IF EXISTS trigger_notify_new_command ON message_queue;
DROP FUNCTION IF EXISTS notify_new_command();
DROP TABLE IF EXISTS message_queue;
