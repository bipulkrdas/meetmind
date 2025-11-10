-- +migrate Up
ALTER TABLE messages
ADD COLUMN message_type TEXT NOT NULL DEFAULT 'user_message',
ADD COLUMN extra_data JSONB;

-- +migrate Down
ALTER TABLE messages
DROP COLUMN message_type,
DROP COLUMN extra_data;
