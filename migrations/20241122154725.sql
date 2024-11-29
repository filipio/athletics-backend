-- Modify "answers" table
ALTER TABLE "answers" ADD COLUMN "points" bigint NOT NULL DEFAULT 0;
-- Modify "questions" table
ALTER TABLE "questions" ADD COLUMN "points" bigint NOT NULL DEFAULT 1;
