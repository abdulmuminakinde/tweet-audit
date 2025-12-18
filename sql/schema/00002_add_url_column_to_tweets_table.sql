-- +goose Up
-- +goose StatementBegin
ALTER TABLE tweets ADD COLUMN url TEXT NOT NULL;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER table tweets DROP COLUMN url;
-- +goose StatementEnd


