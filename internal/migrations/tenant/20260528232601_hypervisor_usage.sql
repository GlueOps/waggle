-- +goose Up
-- modify "hypervisors" table
ALTER TABLE "hypervisors" ADD COLUMN "cpu_used" bigint NOT NULL DEFAULT 0, ADD COLUMN "ram_gb_used" bigint NOT NULL DEFAULT 0, ADD COLUMN "disk_gb_used" bigint NOT NULL DEFAULT 0;

-- +goose Down
-- reverse: modify "hypervisors" table
ALTER TABLE "hypervisors" DROP COLUMN "disk_gb_used", DROP COLUMN "ram_gb_used", DROP COLUMN "cpu_used";
