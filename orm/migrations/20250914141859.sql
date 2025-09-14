-- Create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "firstname" character varying NOT NULL,
  "lastname" character varying NOT NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "reset_token" text NULL,
  "date_of_birth" character varying NOT NULL,
  "profile_image" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Set comment to column: "reset_token" on table: "users"
COMMENT ON COLUMN "users"."reset_token" IS 'Used for resetting password';
