-- +goose Up
-- create "accounts" table
CREATE TABLE "accounts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "email" character varying(255) NOT NULL,
  "display_name" character varying(255) NULL,
  "password_hash" text NULL,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_accounts_email" to table: "accounts"
CREATE UNIQUE INDEX "idx_accounts_email" ON "accounts" ("email");
-- create "auth_audit_events" table
CREATE TABLE "auth_audit_events" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NULL,
  "user_id" uuid NULL,
  "event" character varying(64) NOT NULL,
  "outcome" character varying(32) NOT NULL,
  "ip_address" character varying(64) NULL,
  "user_agent" text NULL,
  "metadata" jsonb NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_auth_audit_events_event" to table: "auth_audit_events"
CREATE INDEX "idx_auth_audit_events_event" ON "auth_audit_events" ("event");
-- create index "idx_auth_audit_events_organization_id" to table: "auth_audit_events"
CREATE INDEX "idx_auth_audit_events_organization_id" ON "auth_audit_events" ("organization_id");
-- create index "idx_auth_audit_events_user_id" to table: "auth_audit_events"
CREATE INDEX "idx_auth_audit_events_user_id" ON "auth_audit_events" ("user_id");
-- create "organizations" table
CREATE TABLE "organizations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NULL,
  "status" character varying(32) NOT NULL DEFAULT 'pending',
  "connection_string" text NULL,
  "encrypted_tenant_key" text NULL,
  "tenant_key_iv" text NULL,
  "tenant_key_tag" text NULL,
  "metadata" jsonb NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_organizations_slug" to table: "organizations"
CREATE UNIQUE INDEX "idx_organizations_slug" ON "organizations" ("slug");
-- create index "idx_organizations_status" to table: "organizations"
CREATE INDEX "idx_organizations_status" ON "organizations" ("status");
-- create "token_sessions" table
CREATE TABLE "token_sessions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "organization_id" uuid NULL,
  "refresh_token_hash" character varying(128) NOT NULL,
  "user_agent" text NULL,
  "ip_address" character varying(64) NULL,
  "expires_at" timestamptz NULL,
  "revoked_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_token_sessions_expires_at" to table: "token_sessions"
CREATE INDEX "idx_token_sessions_expires_at" ON "token_sessions" ("expires_at");
-- create index "idx_token_sessions_organization_id" to table: "token_sessions"
CREATE INDEX "idx_token_sessions_organization_id" ON "token_sessions" ("organization_id");
-- create index "idx_token_sessions_refresh_token_hash" to table: "token_sessions"
CREATE UNIQUE INDEX "idx_token_sessions_refresh_token_hash" ON "token_sessions" ("refresh_token_hash");
-- create index "idx_token_sessions_user_id" to table: "token_sessions"
CREATE INDEX "idx_token_sessions_user_id" ON "token_sessions" ("user_id");
-- create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NOT NULL,
  "email" character varying(255) NOT NULL,
  "display_name" character varying(255) NULL,
  "password_hash" text NULL,
  "is_active" boolean NULL DEFAULT true,
  "last_login_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_users_email" to table: "users"
CREATE INDEX "idx_users_email" ON "users" ("email");
-- create index "idx_users_organization_id" to table: "users"
CREATE INDEX "idx_users_organization_id" ON "users" ("organization_id");

-- +goose Down
-- reverse: create index "idx_users_organization_id" to table: "users"
DROP INDEX "idx_users_organization_id";
-- reverse: create index "idx_users_email" to table: "users"
DROP INDEX "idx_users_email";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: create index "idx_token_sessions_user_id" to table: "token_sessions"
DROP INDEX "idx_token_sessions_user_id";
-- reverse: create index "idx_token_sessions_refresh_token_hash" to table: "token_sessions"
DROP INDEX "idx_token_sessions_refresh_token_hash";
-- reverse: create index "idx_token_sessions_organization_id" to table: "token_sessions"
DROP INDEX "idx_token_sessions_organization_id";
-- reverse: create index "idx_token_sessions_expires_at" to table: "token_sessions"
DROP INDEX "idx_token_sessions_expires_at";
-- reverse: create "token_sessions" table
DROP TABLE "token_sessions";
-- reverse: create index "idx_organizations_status" to table: "organizations"
DROP INDEX "idx_organizations_status";
-- reverse: create index "idx_organizations_slug" to table: "organizations"
DROP INDEX "idx_organizations_slug";
-- reverse: create "organizations" table
DROP TABLE "organizations";
-- reverse: create index "idx_auth_audit_events_user_id" to table: "auth_audit_events"
DROP INDEX "idx_auth_audit_events_user_id";
-- reverse: create index "idx_auth_audit_events_organization_id" to table: "auth_audit_events"
DROP INDEX "idx_auth_audit_events_organization_id";
-- reverse: create index "idx_auth_audit_events_event" to table: "auth_audit_events"
DROP INDEX "idx_auth_audit_events_event";
-- reverse: create "auth_audit_events" table
DROP TABLE "auth_audit_events";
-- reverse: create index "idx_accounts_email" to table: "accounts"
DROP INDEX "idx_accounts_email";
-- reverse: create "accounts" table
DROP TABLE "accounts";
