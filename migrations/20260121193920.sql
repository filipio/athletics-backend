-- Create "refresh_tokens" table
CREATE TABLE "refresh_tokens" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "token_hash" text NOT NULL,
  "session_id" uuid NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "revoked_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_refresh_tokens_session_id" UNIQUE ("session_id")
);
-- Create index "idx_refresh_tokens_expires_at" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_expires_at" ON "refresh_tokens" ("expires_at");
-- Create index "idx_refresh_tokens_revoked_at" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_revoked_at" ON "refresh_tokens" ("revoked_at");
-- Create index "idx_refresh_tokens_session_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_session_id" ON "refresh_tokens" ("session_id");
-- Create index "idx_refresh_tokens_token_hash" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_token_hash" ON "refresh_tokens" ("token_hash");
-- Create index "idx_refresh_tokens_user_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_user_id" ON "refresh_tokens" ("user_id");
