-- +goose Up
-- +goose StatementBegin
CREATE TABLE song (
    id INTEGER PRIMARY KEY ASC,
    title VARCHAR(255) NOT NULL,
    album_id NUMERIC NOT NULL,
    FOREIGN KEY (album_id) REFERENCES album(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE song;
-- +goose StatementEnd
