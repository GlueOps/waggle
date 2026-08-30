-- +goose Up
-- A placement's "hypervisor_id" is the ASSIGNMENT: the hypervisor Waggle booked
-- for this VM, and the value the provisioning pipeline is obliged to build on.
-- It is part of the pipeline's own Terraform state, so it must never be
-- rewritten server-side.
--
-- "observed_hypervisor_id" is where discovery last actually FOUND the guest.
-- NULL means "on its assignment, or not yet seen". A non-NULL value that
-- differs from hypervisor_id is a violated booking: something built the guest
-- in the wrong place or moved it afterwards. Capacity is charged against the
-- observed hypervisor so the scheduler stops overselling the host that is
-- really carrying the guest, while the assignment stays intact as the record of
-- what was promised.
-- modify "placements" table
ALTER TABLE "placements" ADD COLUMN "observed_hypervisor_id" uuid NULL;
-- create index "idx_placements_observed_hypervisor_id" to table: "placements"
CREATE INDEX "idx_placements_observed_hypervisor_id" ON "placements" ("observed_hypervisor_id");

-- +goose Down
-- reverse: modify "placements" table
DROP INDEX "idx_placements_observed_hypervisor_id";
ALTER TABLE "placements" DROP COLUMN "observed_hypervisor_id";
