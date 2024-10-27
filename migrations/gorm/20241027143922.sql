-- Create "events" table
CREATE TABLE "events" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "deadline" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
