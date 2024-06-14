-- +goose Up
-- +goose StatementBegin
CREATE TABLE cache (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP,
    cache_key VARCHAR(1024),
    cache_value TEXT
);
CREATE INDEX idx_cache_key ON cache (cache_key);
CREATE INDEX idx_cache_valid_until ON cache (valid_until);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cache;
-- +goose StatementEnd
