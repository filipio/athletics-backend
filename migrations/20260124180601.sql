-- Modify "events" table
ALTER TABLE "events" ADD COLUMN "status" text NOT NULL DEFAULT 'draft';
