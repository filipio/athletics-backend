-- Add "status" column to "events" table
ALTER TABLE "events" ADD COLUMN "status" text NOT NULL DEFAULT 'draft';
-- Add check constraint for valid status values
ALTER TABLE "events" ADD CONSTRAINT check_event_status CHECK (status IN ('draft', 'published'));
