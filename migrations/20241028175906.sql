-- Create "questions" table
CREATE TABLE "questions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "event_id" bigint NOT NULL,
  "content" text NOT NULL,
  "correct_answer" jsonb NULL,
  "type" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_events_questions" FOREIGN KEY ("event_id") REFERENCES "events" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
