-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories ADD COLUMN is_active BOOLEAN DEFAULT true;
ALTER TABLE products RENAME COLUMN is_available TO is_active;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products RENAME COLUMN is_active TO is_available;
ALTER TABLE categories DROP COLUMN is_active;
-- +goose StatementEnd
