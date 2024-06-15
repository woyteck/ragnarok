-- +goose Up
-- +goose StatementBegin
CREATE TYPE memory_type AS ENUM ('web_article', 'text_file');

CREATE TABLE memories (
    uuid UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    memory_type memory_type NOT NULL,
    source varchar(1024),
    content TEXT
);

CREATE TABLE memory_fragments (
    uuid UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    content_original TEXT,
    content_refined TEXT,
    is_refined BOOLEAN DEFAULT false,
    is_embedded BOOLEAN DEFAULT false,
    memory_id UUID REFERENCES memories(uuid) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE memory_fragments;
DROP TABLE memories;
DROP TYPE memory_type
-- +goose StatementEnd
