CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar(50) UNIQUE NOT NULL,
  "first_name" varchar(50) NOT NULL,
  "last_name" varchar(50) NOT NULL,
  "gender" varchar(50) NOT NULL,
  "birth_date" date NOT NULL,
  "country" varchar(50) NOT NULL,
  "is_mr_user" smallint NOT NULL,
  "is_mobile_user" smallint NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "goals" varchar(500) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "users_profile_mr" (
  "user_id" bigint UNIQUE NOT NULL,
  "total_time_spent_in_mins" bigint NOT NULL DEFAULT 0,
  "deleted_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "users_profile_mobile" (
  "user_id" bigint UNIQUE NOT NULL,
  "total_time_spent_in_mins" bigint NOT NULL DEFAULT 0,
  "deleted_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "session_logs" (
  "uuid" UUID PRIMARY KEY NOT NULL,
  "user_id" bigint NOT NULL,
  "session_type" varchar(50) NOT NULL,
  "session_platform" varchar(20) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "tibetan_singing_bowl_mr" (
  "unique_id" UUID PRIMARY KEY NOT NULL,
  "uuid" UUID UNIQUE NOT NULL,
  "start_mood_rating" smallint NOT NULL,
  "start_mood" varchar(20) NOT NULL,
  "finish_mood_rating" smallint NOT NULL,
  "finish_mood" varchar(20) NOT NULL,
  "session_completed" smallint NOT NULL DEFAULT 0,
  "started_at" timestamptz NOT NULL DEFAULT (now()),
  "ends_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "deleted_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "tummo_breathing_mr" (
  "unique_id" UUID PRIMARY KEY NOT NULL,
  "uuid" UUID UNIQUE NOT NULL,
  "start_mood_rating" smallint NOT NULL,
  "start_mood" varchar(20) NOT NULL,
  "finish_mood_rating" smallint NOT NULL,
  "finish_mood" varchar(20) NOT NULL,
  "session_completed" smallint NOT NULL DEFAULT 0,
  "started_at" timestamptz NOT NULL DEFAULT (now()),
  "ends_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "deleted_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "email_registrations" (
  "unique_id" UUID PRIMARY KEY NOT NULL,
  "email" VARCHAR(50) NOT NULL UNIQUE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "deleted_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "chakra_test_results" (
    "unique_id" UUID PRIMARY KEY NOT NULL,
    "email" VARCHAR(50) NOT NULL,
    "chakra_name" VARCHAR(50) NOT NULL,
    "chakra_score" INTEGER NOT NULL DEFAULT 0,
    "chakra_status" VARCHAR(20) NOT NULL DEFAULT 'inactive',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("country");

CREATE INDEX ON "users" ("is_mr_user");

CREATE INDEX ON "users" ("is_mobile_user");

CREATE INDEX ON "users" ("gender");

CREATE INDEX ON "users" ("birth_date");

CREATE UNIQUE INDEX ON "session_logs" ("uuid", "session_type");

CREATE INDEX ON "session_logs" ("session_type");

CREATE INDEX ON "session_logs" ("session_platform");

CREATE INDEX ON "tibetan_singing_bowl_mr" ("uuid");

CREATE INDEX ON "tibetan_singing_bowl_mr" ("session_completed");

CREATE INDEX ON "tibetan_singing_bowl_mr" ("finish_mood_rating");

CREATE INDEX ON "tummo_breathing_mr" ("uuid");

CREATE INDEX ON "tummo_breathing_mr" ("session_completed");

CREATE INDEX ON "tummo_breathing_mr" ("finish_mood_rating");

COMMENT ON COLUMN "users"."is_mr_user" IS '1 = yes. 0 = no';

COMMENT ON COLUMN "users"."is_mobile_user" IS '1 = yes. 0 = no';

COMMENT ON COLUMN "session_logs"."session_type" IS 'references the name of the table of the session';

COMMENT ON COLUMN "session_logs"."session_platform" IS 'either mr or mobile';

COMMENT ON COLUMN "tibetan_singing_bowl_mr"."start_mood_rating" IS 'between 1-10';

COMMENT ON COLUMN "tibetan_singing_bowl_mr"."finish_mood_rating" IS 'between 1-10';

COMMENT ON COLUMN "tibetan_singing_bowl_mr"."session_completed" IS '0 = incomplete. 1 = complete';

COMMENT ON COLUMN "tummo_breathing_mr"."start_mood_rating" IS 'between 1-10';

COMMENT ON COLUMN "tummo_breathing_mr"."finish_mood_rating" IS 'between 1-10';

COMMENT ON COLUMN "tummo_breathing_mr"."session_completed" IS '0 = incomplete. 1 = complete';

ALTER TABLE "users_profile_mr" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "users_profile_mobile" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "session_logs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "tibetan_singing_bowl_mr" ADD FOREIGN KEY ("uuid") REFERENCES "session_logs" ("uuid") ON DELETE CASCADE;

ALTER TABLE "tummo_breathing_mr" ADD FOREIGN KEY ("uuid") REFERENCES "session_logs" ("uuid") ON DELETE CASCADE;


-- ALTER TABLE "session_logs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

-- ALTER TABLE "tibetan_singing_bowl_mr" ADD FOREIGN KEY ("uuid") REFERENCES "session_logs" ("uuid") ON DELETE CASCADE;

-- ALTER TABLE "tummo_breathing_mr" ADD FOREIGN KEY ("uuid") REFERENCES "session_logs" ("uuid") ON DELETE CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE "session_logs"
  ALTER COLUMN "uuid" SET DEFAULT uuid_generate_v4();

ALTER TABLE "tibetan_singing_bowl_mr"
  ALTER COLUMN "unique_id" SET DEFAULT uuid_generate_v4();

ALTER TABLE "tummo_breathing_mr"
  ALTER COLUMN "unique_id" SET DEFAULT uuid_generate_v4();

-- Create the trigger function for soft deletion
CREATE OR REPLACE FUNCTION soft_delete_user()
RETURNS TRIGGER AS $$
BEGIN
    -- Soft delete user
    UPDATE users SET deleted_at = NOW() WHERE id = OLD.id AND deleted_at IS NULL;

    -- Soft delete session logs linked to the user
    UPDATE session_logs 
    SET deleted_at = NOW() 
    WHERE user_id = OLD.id AND deleted_at IS NULL;

    -- Soft delete related users_profile_mr
    UPDATE users_profile_mr 
    SET deleted_at = NOW() 
    WHERE user_id IN (SELECT user_id FROM users_profile_mr WHERE user_id = OLD.id) 
      AND deleted_at IS NULL;

    -- Soft delete related users_profile_mobile
    UPDATE users_profile_mobile
    SET deleted_at = NOW() 
    WHERE user_id IN (SELECT user_id FROM users_profile_mobile WHERE user_id = OLD.id) 
      AND deleted_at IS NULL;

    -- Soft delete related tibetan singing bowl entries
    UPDATE tibetan_singing_bowl_mr 
    SET deleted_at = NOW() 
    WHERE uuid IN (SELECT uuid FROM session_logs WHERE user_id = OLD.id) 
      AND deleted_at IS NULL;

    -- Soft delete related tummo breathing entries
    UPDATE tummo_breathing_mr 
    SET deleted_at = NOW() 
    WHERE uuid IN (SELECT uuid FROM session_logs WHERE user_id = OLD.id) 
      AND deleted_at IS NULL;

    RETURN NULL;  -- Indicate that the default action should not be taken
END;
$$ LANGUAGE plpgsql;

-- Create the trigger on the users table
CREATE TRIGGER trigger_soft_delete_users
BEFORE DELETE ON users
FOR EACH ROW EXECUTE FUNCTION soft_delete_user();


-- Create the trigger function for soft deletion of session logs and related sessions data
CREATE OR REPLACE FUNCTION soft_delete_session_logs_and_related_sessions_data()
RETURNS TRIGGER AS $$
BEGIN
    -- Soft delete session_logs
    UPDATE session_logs SET deleted_at = NOW() WHERE uuid = OLD.uuid AND deleted_at IS NULL;
    -- Soft delete related tibetan singing bowl entries
    UPDATE tibetan_singing_bowl_mr 
    SET deleted_at = NOW() 
    WHERE uuid = OLD.uuid 
      AND deleted_at IS NULL;

    -- Soft delete related tummo breathing entries
    UPDATE tummo_breathing_mr 
    SET deleted_at = NOW() 
    WHERE uuid = OLD.uuid 
      AND deleted_at IS NULL;

    RETURN NULL;  -- Indicate that the default action should not be taken
END;
$$ LANGUAGE plpgsql;

-- Create the trigger on the session_logs table
CREATE TRIGGER trigger_soft_delete_session_logs
BEFORE DELETE ON session_logs
FOR EACH ROW EXECUTE FUNCTION soft_delete_session_logs_and_related_sessions_data();


-- -- users table - Soft deletion for users
-- CREATE RULE users_soft_deletion AS ON DELETE TO users
-- DO INSTEAD (
--   UPDATE users 
--   SET deleted_at = NOW() 
--   WHERE id = OLD.id AND deleted_at IS NULL

--   execute session_logs_soft_deletion if not executed
--    then execute session_logs_soft_deletion_for_tibetan_singing_bowl_mr
--    then execute 
-- );


-- -- users table -  Cascade soft deletion for session_logs
-- CREATE RULE users_session_logs_soft_deletion_for_session_logs AS ON DELETE TO users
-- DO ALSO (
--   UPDATE session_logs 
--   SET deleted_at = NOW() 
--   WHERE user_id = OLD.id AND deleted_at IS NULL
-- );


-- -- session_logs - Soft deletion for session_logs
-- CREATE RULE session_logs_soft_deletion AS ON DELETE TO session_logs
-- DO INSTEAD (
--   UPDATE session_logs 
--   SET deleted_at = NOW() 
--   WHERE uuid = OLD.uuid AND deleted_at IS NULL
-- );

-- -- session_logs - Cascade soft deletion for tibetan_singing_bowl_mr
-- CREATE RULE session_logs_soft_deletion_for_tibetan_singing_bowl_mr AS ON DELETE TO session_logs
-- DO ALSO (
--   UPDATE tibetan_singing_bowl_mr 
--   SET deleted_at = NOW() 
--   WHERE uuid = OLD.uuid AND deleted_at IS NULL
-- );

-- -- session_logs - Cascade soft deletion for tummo_breathing_mr
-- CREATE RULE session_logs_soft_deletion_for_tummo_breathing_mr AS ON DELETE TO session_logs
-- DO ALSO (
--   UPDATE tummo_breathing_mr 
--   SET deleted_at = NOW() 
--   WHERE uuid = OLD.uuid AND deleted_at IS NULL
-- );


-- -- tibetan_singing_bowl_mr - Soft deletion for tibetan_singing_bowl_mr
-- CREATE RULE tibetan_singing_bowl_mr_soft_deletion AS ON DELETE TO tibetan_singing_bowl_mr
-- DO INSTEAD (
--   UPDATE tibetan_singing_bowl_mr 
--   SET deleted_at = NOW() 
--   WHERE uuid = OLD.uuid AND deleted_at IS NULL
-- );

-- -- tummo_breathing_mr - Soft deletion for tummo_breathing_mr
-- CREATE RULE tummo_breathing_mr_soft_deletion AS ON DELETE TO tummo_breathing_mr
-- DO INSTEAD (
--   UPDATE tummo_breathing_mr 
--   SET deleted_at = NOW() 
--   WHERE uuid = OLD.uuid AND deleted_at IS NULL
-- );