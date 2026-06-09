-- +goose Up
-- modify "hypervisors" table
ALTER TABLE "hypervisors" ADD COLUMN "schedulable" boolean NOT NULL DEFAULT true;

-- +goose Down
-- reverse: modify "hypervisors" table
ALTER TABLE "hypervisors" DROP COLUMN "schedulable";
