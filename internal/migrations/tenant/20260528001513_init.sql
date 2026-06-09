-- +goose Up
-- create "datacenters" table
CREATE TABLE "datacenters" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NULL,
  "url" character varying(255) NULL,
  "encrypted_token_key" text NULL,
  "token_key_iv" text NULL,
  "token_key_tag" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "hypervisors" table
CREATE TABLE "hypervisors" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "datacenter_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "cpu_total" bigint NOT NULL DEFAULT 0,
  "cpu_reserved" bigint NOT NULL DEFAULT 0,
  "ram_gb_total" bigint NOT NULL DEFAULT 0,
  "ram_gb_reserved" bigint NOT NULL DEFAULT 0,
  "disk_gb_total" bigint NOT NULL DEFAULT 0,
  "disk_gb_reserved" bigint NOT NULL DEFAULT 0,
  "last_synced_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "ux_hypervisor_datacenter_name" to table: "hypervisors"
CREATE UNIQUE INDEX "ux_hypervisor_datacenter_name" ON "hypervisors" ("datacenter_id", "name");
-- create "placements" table
CREATE TABLE "placements" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "pool_id" uuid NOT NULL,
  "hypervisor_id" uuid NOT NULL,
  "vmid" bigint NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_placements_hypervisor_id" to table: "placements"
CREATE INDEX "idx_placements_hypervisor_id" ON "placements" ("hypervisor_id");
-- create index "idx_placements_pool_id" to table: "placements"
CREATE INDEX "idx_placements_pool_id" ON "placements" ("pool_id");
-- create index "idx_placements_vm_id" to table: "placements"
CREATE INDEX "idx_placements_vm_id" ON "placements" ("vmid");
-- create "pools" table
CREATE TABLE "pools" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "datacenter_id" uuid NOT NULL,
  "slot_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "desired_count" bigint NOT NULL DEFAULT 0,
  "metadata" jsonb NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_pools_datacenter_id" to table: "pools"
CREATE INDEX "idx_pools_datacenter_id" ON "pools" ("datacenter_id");
-- create index "idx_pools_name" to table: "pools"
CREATE INDEX "idx_pools_name" ON "pools" ("name");
-- create index "idx_pools_slot_id" to table: "pools"
CREATE INDEX "idx_pools_slot_id" ON "pools" ("slot_id");
-- create "slots" table
CREATE TABLE "slots" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "vcpu" bigint NOT NULL,
  "ram_gb" bigint NOT NULL,
  "disk_gb" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_slots_name" to table: "slots"
CREATE UNIQUE INDEX "idx_slots_name" ON "slots" ("name");

-- +goose Down
-- reverse: create index "idx_slots_name" to table: "slots"
DROP INDEX "idx_slots_name";
-- reverse: create "slots" table
DROP TABLE "slots";
-- reverse: create index "idx_pools_slot_id" to table: "pools"
DROP INDEX "idx_pools_slot_id";
-- reverse: create index "idx_pools_name" to table: "pools"
DROP INDEX "idx_pools_name";
-- reverse: create index "idx_pools_datacenter_id" to table: "pools"
DROP INDEX "idx_pools_datacenter_id";
-- reverse: create "pools" table
DROP TABLE "pools";
-- reverse: create index "idx_placements_vm_id" to table: "placements"
DROP INDEX "idx_placements_vm_id";
-- reverse: create index "idx_placements_pool_id" to table: "placements"
DROP INDEX "idx_placements_pool_id";
-- reverse: create index "idx_placements_hypervisor_id" to table: "placements"
DROP INDEX "idx_placements_hypervisor_id";
-- reverse: create "placements" table
DROP TABLE "placements";
-- reverse: create index "ux_hypervisor_datacenter_name" to table: "hypervisors"
DROP INDEX "ux_hypervisor_datacenter_name";
-- reverse: create "hypervisors" table
DROP TABLE "hypervisors";
-- reverse: create "datacenters" table
DROP TABLE "datacenters";
