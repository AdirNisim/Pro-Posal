-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories ADD COLUMN type TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories DROP COLUMN type;
-- +goose StatementEnd
