-- +goose Up
-- +goose StatementBegin
CREATE TABLE album (
    id INTEGER PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    year NUMERIC NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE album;
-- +goose StatementEnd
