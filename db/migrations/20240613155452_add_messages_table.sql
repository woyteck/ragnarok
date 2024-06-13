-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    uuid UUID PRIMARY KEY,
    conversation_id UUID REFERENCES conversations(uuid),
    role VARCHAR(50),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
