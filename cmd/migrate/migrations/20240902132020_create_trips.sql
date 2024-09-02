-- +goose Up
-- +goose StatementBegin
CREATE TABLE "trips" (
    "id" BIGSERIAL PRIMARY KEY,
    "trip_code" TEXT NOT NULL,
    "user_id" bigint NOT NULL,
    "driver_id" bigint,
    "status" VARCHAR(20) NOT NULL,
    "pickup_location" geometry NOT NULL,
    "pickup_address" TEXT NOT NULL,
    "dropoff_location" geometry NOT NULL,
    "dropoff_address" TEXT NOT NULL,
    "request_time" timestamp(0) NOT NULL DEFAULT now(),
    "start_time" timestamp(0),
    "end_time" timestamp(0),
    "price" DECIMAL(10, 2),
    "distance" DECIMAL(10, 2),
    "created_at" timestamp(0) NOT NULL DEFAULT now(),
    "updated_at" timestamp(0) NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX unique_trips_trip_code_idx ON trips (trip_code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "trips";
-- +goose StatementEnd
