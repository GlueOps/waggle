-- +goose Up
-- modify "users" table
ALTER TABLE "users" ADD COLUMN "role" character varying(16) NOT NULL DEFAULT 'member';
-- backfill: every pre-existing membership was created by signup (the org
-- creator), so promote them to owner rather than the 'member' default.
UPDATE "users" SET "role" = 'owner';

-- +goose Down
-- reverse: modify "users" table
ALTER TABLE "users" DROP COLUMN "role";
