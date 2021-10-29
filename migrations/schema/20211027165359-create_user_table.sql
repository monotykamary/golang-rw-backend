
-- +migrate Up
CREATE TABLE "users" (
  "id" uuid NOT NULL PRIMARY KEY DEFAULT md5(random()::text || clock_timestamp()::text)::uuid,
  "email" VARCHAR NOT NULL,
  "hashed_password" VARCHAR NOT NULL,
  "status" VARCHAR DEFAULT 'inactive',
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);


-- +migrate Down

DROP TABLE IF EXISTS "users";
