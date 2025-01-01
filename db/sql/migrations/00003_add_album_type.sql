-- +goose Up
-- +goose StatementBegin
CREATE TABLE type (
    id INTEGER PRIMARY KEY,
    name VARCHAR(25) NOT NULL
);
ALTER TABLE album ADD COLUMN type_id INTEGER REFERENCES type(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE album DROP COLUMN type_id;
DROP TABLE type;
-- +goose StatementEnd
