-- +goose Up
-- modify "accounts" table
ALTER TABLE "accounts" DROP COLUMN "email", ADD COLUMN "last_login_at" timestamptz NULL;
-- modify "organizations" table
ALTER TABLE "organizations" ADD COLUMN "domain" character varying(255) NULL;
-- create index "idx_organizations_domain" to table: "organizations"
CREATE UNIQUE INDEX "idx_organizations_domain" ON "organizations" ("domain");
-- modify "token_sessions" table
ALTER TABLE "token_sessions" DROP COLUMN "user_id", ADD COLUMN "account_id" uuid NOT NULL;
-- create index "idx_token_sessions_account_id" to table: "token_sessions"
CREATE INDEX "idx_token_sessions_account_id" ON "token_sessions" ("account_id");
-- modify "users" table
ALTER TABLE "users" DROP COLUMN "email", DROP COLUMN "password_hash", ADD COLUMN "account_id" uuid NOT NULL;
-- create index "ux_user_account_org" to table: "users"
CREATE UNIQUE INDEX "ux_user_account_org" ON "users" ("account_id", "organization_id");
-- create "account_emails" table
CREATE TABLE "account_emails" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "account_id" uuid NOT NULL,
  "email" character varying(255) NOT NULL,
  "is_primary" boolean NOT NULL DEFAULT false,
  "verified_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_account_emails_account_id" to table: "account_emails"
CREATE INDEX "idx_account_emails_account_id" ON "account_emails" ("account_id");
-- create index "idx_account_emails_email" to table: "account_emails"
CREATE UNIQUE INDEX "idx_account_emails_email" ON "account_emails" ("email");
-- create index "ux_account_emails_primary" to table: "account_emails"
CREATE UNIQUE INDEX "ux_account_emails_primary" ON "account_emails" ("account_id") WHERE is_primary;

-- +goose Down
-- reverse: create index "ux_account_emails_primary" to table: "account_emails"
DROP INDEX "ux_account_emails_primary";
-- reverse: create index "idx_account_emails_email" to table: "account_emails"
DROP INDEX "idx_account_emails_email";
-- reverse: create index "idx_account_emails_account_id" to table: "account_emails"
DROP INDEX "idx_account_emails_account_id";
-- reverse: create "account_emails" table
DROP TABLE "account_emails";
-- reverse: create index "ux_user_account_org" to table: "users"
DROP INDEX "ux_user_account_org";
-- reverse: modify "users" table
ALTER TABLE "users" DROP COLUMN "account_id", ADD COLUMN "password_hash" text NULL, ADD COLUMN "email" character varying(255) NOT NULL;
-- reverse: create index "idx_token_sessions_account_id" to table: "token_sessions"
DROP INDEX "idx_token_sessions_account_id";
-- reverse: modify "token_sessions" table
ALTER TABLE "token_sessions" DROP COLUMN "account_id", ADD COLUMN "user_id" uuid NOT NULL;
-- reverse: create index "idx_organizations_domain" to table: "organizations"
DROP INDEX "idx_organizations_domain";
-- reverse: modify "organizations" table
ALTER TABLE "organizations" DROP COLUMN "domain";
-- reverse: modify "accounts" table
ALTER TABLE "accounts" DROP COLUMN "last_login_at", ADD COLUMN "email" character varying(255) NOT NULL;
