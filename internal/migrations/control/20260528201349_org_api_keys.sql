-- +goose Up
-- create "org_api_keys" table
CREATE TABLE "org_api_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "token_hash" character varying(64) NOT NULL,
  "prefix" character varying(16) NOT NULL,
  "created_by_account_id" uuid NULL,
  "last_used_at" timestamptz NULL,
  "expires_at" timestamptz NULL,
  "revoked_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_org_api_keys_expires_at" to table: "org_api_keys"
CREATE INDEX "idx_org_api_keys_expires_at" ON "org_api_keys" ("expires_at");
-- create index "idx_org_api_keys_organization_id" to table: "org_api_keys"
CREATE INDEX "idx_org_api_keys_organization_id" ON "org_api_keys" ("organization_id");
-- create index "idx_org_api_keys_token_hash" to table: "org_api_keys"
CREATE UNIQUE INDEX "idx_org_api_keys_token_hash" ON "org_api_keys" ("token_hash");

-- +goose Down
-- reverse: create index "idx_org_api_keys_token_hash" to table: "org_api_keys"
DROP INDEX "idx_org_api_keys_token_hash";
-- reverse: create index "idx_org_api_keys_organization_id" to table: "org_api_keys"
DROP INDEX "idx_org_api_keys_organization_id";
-- reverse: create index "idx_org_api_keys_expires_at" to table: "org_api_keys"
DROP INDEX "idx_org_api_keys_expires_at";
-- reverse: create "org_api_keys" table
DROP TABLE "org_api_keys";
