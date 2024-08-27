-- +goose Up
-- +goose StatementBegin
CREATE TABLE "locations" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "geometry" geometry NOT NULL,
    "created_at" timestamp(0) NOT NULL DEFAULT now(),
    "updated_at" timestamp(0) NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX unique_user_id_idx ON locations (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "locations";
-- +goose StatementEnd
