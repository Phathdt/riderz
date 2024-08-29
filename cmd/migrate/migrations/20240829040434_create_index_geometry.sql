-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_locations_geometry ON locations USING GIST (geometry);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "public"."idx_locations_geometry";
-- +goose StatementEnd
