-- Create "tasks" table
CREATE TABLE "tasks" (
  "id" text NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "title" text NULL,
  "description" text NULL,
  "status" bigint NULL,
  "recurrence" bigint NULL,
  "scheduled_time" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_tasks_status" to table: "tasks"
CREATE INDEX "idx_tasks_status" ON "tasks" ("status");
