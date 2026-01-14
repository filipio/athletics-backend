-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "email_verified" boolean NOT NULL DEFAULT false;
-- Create "email_verification_rate_limits" table
CREATE TABLE "email_verification_rate_limits" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "email" text NOT NULL,
  "request_count" bigint NOT NULL DEFAULT 0,
  "last_request_at" timestamptz NULL,
  "blocked_until" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_email_verification_rate_limits_email" UNIQUE ("email")
);
-- Create index "idx_email_verification_rate_limits_blocked_until" to table: "email_verification_rate_limits"
CREATE INDEX "idx_email_verification_rate_limits_blocked_until" ON "email_verification_rate_limits" ("blocked_until");
-- Create index "idx_email_verification_rate_limits_email" to table: "email_verification_rate_limits"
CREATE INDEX "idx_email_verification_rate_limits_email" ON "email_verification_rate_limits" ("email");
-- Create "pending_registrations" table
CREATE TABLE "pending_registrations" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "email" text NOT NULL,
  "username" text NOT NULL,
  "password_hash" text NOT NULL,
  "verification_token" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "verified" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- Create index "idx_pending_registrations_email" to table: "pending_registrations"
CREATE INDEX "idx_pending_registrations_email" ON "pending_registrations" ("email");
-- Create index "idx_pending_registrations_expires_at" to table: "pending_registrations"
CREATE INDEX "idx_pending_registrations_expires_at" ON "pending_registrations" ("expires_at");
-- Create index "idx_pending_registrations_verification_token" to table: "pending_registrations"
CREATE INDEX "idx_pending_registrations_verification_token" ON "pending_registrations" ("verification_token");
