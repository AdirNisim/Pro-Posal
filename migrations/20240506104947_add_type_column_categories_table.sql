-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories ADD COLUMN type VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories DROP COLUMN type;
-- +goose StatementEnd
