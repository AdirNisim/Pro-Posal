-- +goose Up
-- +goose StatementBegin
CREATE TABLE "session"(
    "id" UUID NOT NULL PRIMARY KEY,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "expires_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "session";
-- +goose StatementEnd
