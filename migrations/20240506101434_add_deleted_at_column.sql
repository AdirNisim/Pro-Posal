-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP(0);
ALTER TABLE companies ADD COLUMN deleted_at TIMESTAMP(0);
ALTER TABLE contract_templates ADD COLUMN deleted_at TIMESTAMP(0);
ALTER TABLE offers ADD COLUMN deleted_at TIMESTAMP(0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN deleted_at;
ALTER TABLE companies DROP COLUMN deleted_at;
ALTER TABLE contract_templates DROP COLUMN deleted_at;
ALTER TABLE offers DROP COLUMN deleted_at;
-- +goose StatementEnd
