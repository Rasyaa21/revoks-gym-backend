-- Migration script for Revok's Gym database schema
-- Run this file before starting the application to create PostgreSQL enums and constraints

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =========================
-- ENUMS
-- =========================
DO $$ BEGIN CREATE TYPE user_role_enum AS ENUM ('SUPERADMIN','ADMIN','CASHIER','MEMBER','PT'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE user_status_enum AS ENUM ('active','blocked'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE membership_status_enum AS ENUM ('ACTIVE','EXPIRED','CANCELLED'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE order_status_enum AS ENUM ('PENDING','PAID','CANCELLED','EXPIRED'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE payment_status_enum AS ENUM ('PENDING','PAID','FAILED'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE attendance_direction_enum AS ENUM ('IN','OUT'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE attendance_result_enum AS ENUM ('ACCEPTED','REJECTED'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE otv_pass_status_enum AS ENUM ('UNUSED','USED','EXPIRED'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Note: Tables will be created automatically by GORM AutoMigrate
-- However, we need to add custom constraints after GORM creates the tables

-- Add constraints after running the application for the first time:
-- ALTER TABLE attendance_logs DROP CONSTRAINT IF EXISTS chk_attendance_actor;
-- ALTER TABLE attendance_logs ADD CONSTRAINT chk_attendance_actor
-- CHECK (
--   (user_id IS NOT NULL AND one_time_pass_id IS NULL)
--   OR
--   (user_id IS NULL AND one_time_pass_id IS NOT NULL)
-- );

-- ALTER TABLE attendance_logs DROP CONSTRAINT IF EXISTS chk_attendance_direction;
-- ALTER TABLE attendance_logs ADD CONSTRAINT chk_attendance_direction
-- CHECK (
--   (user_id IS NOT NULL AND direction IS NOT NULL)
--   OR
--   (one_time_pass_id IS NOT NULL AND direction IS NULL)
-- );
