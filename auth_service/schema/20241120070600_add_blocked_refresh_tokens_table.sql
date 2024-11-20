-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS blocked_refresh_tokens(
    token VARCHAR(512) UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blocked_refresh_tokens;
-- +goose StatementEnd
