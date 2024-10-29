-- Create "answers" table
CREATE TABLE "answers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "question_id" bigint NOT NULL,
  "content" jsonb NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_questions_answers" FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_answers" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
