-- +goose Up
-- modify "datacenters" table
ALTER TABLE "datacenters" ADD COLUMN "cpu_overcommit_ratio" numeric(5,2) NOT NULL DEFAULT 1;
-- modify "hypervisors" table
ALTER TABLE "hypervisors" ADD COLUMN "cpu_overcommit_ratio" numeric(5,2) NOT NULL DEFAULT 1;

-- +goose Down
-- reverse: modify "hypervisors" table
ALTER TABLE "hypervisors" DROP COLUMN "cpu_overcommit_ratio";
-- reverse: modify "datacenters" table
ALTER TABLE "datacenters" DROP COLUMN "cpu_overcommit_ratio";
