-- +goose Up
-- +goose StatementBegin
ALTER TABLE chats ALTER COLUMN name DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chats ALTER COLUMN name SET NOT NULL;
-- +goose StatementEnd
