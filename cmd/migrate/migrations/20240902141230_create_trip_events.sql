-- +goose Up
-- +goose StatementBegin
CREATE TABLE "trip_events" (
    "id" BIGSERIAL PRIMARY KEY,
    "trip_id" bigint NOT NULL,
    "event_type" VARCHAR(50) NOT NULL,
    "status" VARCHAR(20),
    "event_data" JSONB,
    "created_at" timestamp(0) NOT NULL DEFAULT now(),
    "updated_at" timestamp(0) NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "trip_events";
-- +goose StatementEnd
