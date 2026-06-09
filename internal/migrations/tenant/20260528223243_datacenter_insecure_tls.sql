-- +goose Up
-- modify "datacenters" table
ALTER TABLE "datacenters" ADD COLUMN "insecure_skip_verify" boolean NOT NULL DEFAULT false;

-- +goose Down
-- reverse: modify "datacenters" table
ALTER TABLE "datacenters" DROP COLUMN "insecure_skip_verify";
