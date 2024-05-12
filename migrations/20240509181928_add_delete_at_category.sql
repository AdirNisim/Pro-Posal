-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories ADD COLUMN deleted_at TIMESTAMP(0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories DROP COLUMN deleted_at;
-- +goose StatementEnd
