-- +goose Up
-- +goose StatementBegin
CREATE TABLE conversations (
    uuid UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE INDEX idx_uuid ON conversations (uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE conversations;
-- +goose StatementEnd
