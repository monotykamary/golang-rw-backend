
-- +migrate Up
CREATE TABLE "booking" (
  "id" uuid NOT NULL PRIMARY KEY DEFAULT md5(random()::text || clock_timestamp()::text)::uuid,
  "user_id" uuid NOT NULL,
  "status" VARCHAR DEFAULT 'booking',
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

ALTER TABLE booking
ADD CONSTRAINT booking_user_fk
FOREIGN KEY (user_id)
REFERENCES users(id);


-- +migrate Down

ALTER TABLE booking
DROP CONSTRAINT KEY booking_user_fk

DROP TABLE IF EXISTS "booking";
