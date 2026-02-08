-- Complete database schema for Revok's Gym
-- This is the full schema from your requirements
-- GORM will auto-migrate most tables, but enums and constraints need to be created manually

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

-- =========================
-- USERS
-- =========================
CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  role user_role_enum NOT NULL,
  email VARCHAR UNIQUE,
  phone VARCHAR UNIQUE,
  password_hash VARCHAR NOT NULL,
  full_name VARCHAR NOT NULL,
  photo_url VARCHAR,
  status user_status_enum NOT NULL DEFAULT 'active',
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- =========================
-- MEMBERS
-- =========================
CREATE TABLE IF NOT EXISTS members (
  member_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID UNIQUE NOT NULL REFERENCES users(user_id),
  legacy_member_no BIGINT,
  legacy_record_id BIGINT,
  gender VARCHAR,
  birth_date DATE,
  address TEXT,
  phone_override VARCHAR,
  qr_static_hash VARCHAR UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- =========================
-- PERSONAL TRAINERS
-- =========================
CREATE TABLE IF NOT EXISTS personal_trainers (
  pt_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID UNIQUE NOT NULL REFERENCES users(user_id),
  bio TEXT,
  experience_years INT,
  specialties TEXT,
  base_rate_amount DECIMAL,
  rating_avg DECIMAL,
  rating_count INT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS pt_reviews (
  review_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  pt_id UUID NOT NULL REFERENCES personal_trainers(pt_id),
  member_id UUID NOT NULL REFERENCES members(member_id),
  rating INT,
  comment TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS pt_availability_rules (
  rule_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  pt_id UUID NOT NULL REFERENCES personal_trainers(pt_id),
  weekday VARCHAR NOT NULL,
  start_time TIME NOT NULL,
  end_time TIME NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- =========================
-- MEMBERSHIP (POS/MOBILE)
-- =========================
CREATE TABLE IF NOT EXISTS membership_plans (
  plan_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR NOT NULL,
  price_amount DECIMAL NOT NULL,
  duration_days INT,
  session_quota INT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- =========================
-- ORDERS (SHARED)
-- =========================
CREATE TABLE IF NOT EXISTS orders (
  order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  channel VARCHAR NOT NULL,
  buyer_user_id UUID REFERENCES users(user_id),
  status order_status_enum NOT NULL DEFAULT 'PENDING',
  total_amount DECIMAL NOT NULL,
  currency VARCHAR NOT NULL DEFAULT 'IDR',
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  expired_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
  order_item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL REFERENCES orders(order_id),
  item_type VARCHAR NOT NULL,
  ref_id UUID NOT NULL,
  qty INT NOT NULL DEFAULT 1,
  unit_price DECIMAL NOT NULL,
  subtotal DECIMAL NOT NULL
);

-- =========================
-- MEMBER MEMBERSHIPS
-- =========================
CREATE TABLE IF NOT EXISTS member_memberships (
  membership_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  member_id UUID NOT NULL REFERENCES members(member_id),
  plan_id UUID NOT NULL REFERENCES membership_plans(plan_id),
  start_at TIMESTAMP,
  expire_at TIMESTAMP,
  status membership_status_enum NOT NULL,
  total_sessions INT,
  used_sessions INT,
  remaining_sessions INT,
  activated_by_order_id UUID REFERENCES orders(order_id),
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS session_usage_logs (
  usage_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  membership_id UUID NOT NULL REFERENCES member_memberships(membership_id),
  used_at TIMESTAMP NOT NULL DEFAULT now(),
  used_by VARCHAR,
  notes VARCHAR
);

-- =========================
-- ONE-TIME VISIT (OTV)
-- =========================
CREATE TABLE IF NOT EXISTS one_time_products (
  otv_product_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR NOT NULL,
  price_amount DECIMAL NOT NULL,
  validity_minutes INT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS one_time_passes (
  one_time_pass_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL REFERENCES orders(order_id),
  otv_product_id UUID NOT NULL REFERENCES one_time_products(otv_product_id),
  token_hash VARCHAR UNIQUE NOT NULL,
  status otv_pass_status_enum NOT NULL DEFAULT 'UNUSED',
  issued_at TIMESTAMP NOT NULL DEFAULT now(),
  expires_at TIMESTAMP,
  used_at TIMESTAMP
);

-- =========================
-- ACCESS / GATE
-- =========================
CREATE TABLE IF NOT EXISTS scanner_devices (
  scanner_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR NOT NULL,
  location VARCHAR,
  api_key_hash VARCHAR NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_access_state (
  user_id UUID PRIMARY KEY REFERENCES users(user_id),
  gate_state VARCHAR NOT NULL DEFAULT 'OUTSIDE',
  last_scanned_at TIMESTAMP,
  last_scanner_id UUID REFERENCES scanner_devices(scanner_id)
);

CREATE TABLE IF NOT EXISTS attendance_logs (
  attendance_log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  scanned_at TIMESTAMP NOT NULL DEFAULT now(),
  scanner_id UUID NOT NULL REFERENCES scanner_devices(scanner_id),
  user_id UUID REFERENCES users(user_id),
  one_time_pass_id UUID REFERENCES one_time_passes(one_time_pass_id),
  scanned_token_hash VARCHAR,
  direction attendance_direction_enum,
  result attendance_result_enum NOT NULL,
  reject_reason VARCHAR
);

ALTER TABLE attendance_logs DROP CONSTRAINT IF EXISTS chk_attendance_actor;
ALTER TABLE attendance_logs ADD CONSTRAINT chk_attendance_actor
CHECK (
  (user_id IS NOT NULL AND one_time_pass_id IS NULL)
  OR
  (user_id IS NULL AND one_time_pass_id IS NOT NULL)
);

ALTER TABLE attendance_logs DROP CONSTRAINT IF EXISTS chk_attendance_direction;
ALTER TABLE attendance_logs ADD CONSTRAINT chk_attendance_direction
CHECK (
  (user_id IS NOT NULL AND direction IS NOT NULL)
  OR
  (one_time_pass_id IS NOT NULL AND direction IS NULL)
);

-- =========================
-- PAYMENTS (MIDTRANS-READY)
-- =========================
CREATE TABLE IF NOT EXISTS payments (
  payment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL REFERENCES orders(order_id),
  provider VARCHAR NOT NULL DEFAULT 'MIDTRANS',
  midtrans_order_id VARCHAR NOT NULL,
  midtrans_transaction_id VARCHAR,
  payment_type VARCHAR,
  transaction_status VARCHAR,
  fraud_status VARCHAR,
  gross_amount DECIMAL NOT NULL,
  currency VARCHAR NOT NULL DEFAULT 'IDR',
  actions_json JSONB,
  va_numbers_json JSONB,
  metadata_json JSONB,
  status payment_status_enum NOT NULL DEFAULT 'PENDING',
  paid_at TIMESTAMP,
  expires_at TIMESTAMP,
  raw_response_json JSONB,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT uq_midtrans_order UNIQUE (midtrans_order_id),
  CONSTRAINT uq_midtrans_trx UNIQUE (midtrans_transaction_id)
);

-- =========================
-- PAYMENT WEBHOOK EVENTS
-- =========================
CREATE TABLE IF NOT EXISTS payment_webhook_events (
  event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  provider VARCHAR NOT NULL DEFAULT 'MIDTRANS',
  provider_event_id VARCHAR UNIQUE,
  order_id UUID REFERENCES orders(order_id),
  midtrans_order_id VARCHAR,
  midtrans_transaction_id VARCHAR,
  raw_payload_json JSONB NOT NULL,
  received_at TIMESTAMP NOT NULL DEFAULT now(),
  processed_at TIMESTAMP,
  process_status VARCHAR NOT NULL DEFAULT 'RECEIVED'
);

-- =========================
-- NOTIFICATIONS
-- =========================
CREATE TABLE IF NOT EXISTS notifications (
  notification_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(user_id),
  type VARCHAR NOT NULL,
  title VARCHAR NOT NULL,
  body TEXT,
  payload_json JSONB,
  read_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- =========================
-- AUDIT LOGS
-- =========================
CREATE TABLE IF NOT EXISTS audit_logs (
  audit_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  actor_user_id UUID REFERENCES users(user_id),
  action VARCHAR NOT NULL,
  entity_type VARCHAR,
  entity_id UUID,
  payload_json JSONB,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- =========================
-- WORKOUT (TEMPLATE-BASED)
-- =========================
CREATE TABLE IF NOT EXISTS exercise_catalog (
  exercise_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR NOT NULL,
  primary_muscle VARCHAR,
  equipment VARCHAR,
  video_url VARCHAR,
  image_url VARCHAR,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS workout_templates (
  template_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR NOT NULL,
  description TEXT,
  focus_muscle VARCHAR,
  level VARCHAR,
  created_by_user_id UUID REFERENCES users(user_id),
  coach_pt_id UUID REFERENCES personal_trainers(pt_id),
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS workout_template_items (
  item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  template_id UUID NOT NULL REFERENCES workout_templates(template_id),
  exercise_id UUID NOT NULL REFERENCES exercise_catalog(exercise_id),
  order_no INT,
  sets INT,
  reps VARCHAR,
  rest_seconds INT,
  tempo VARCHAR,
  rpe INT,
  notes TEXT
);

CREATE TABLE IF NOT EXISTS member_workout_follows (
  follow_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  member_id UUID NOT NULL REFERENCES members(member_id),
  template_id UUID NOT NULL REFERENCES workout_templates(template_id),
  followed_at TIMESTAMP NOT NULL DEFAULT now(),
  status VARCHAR NOT NULL DEFAULT 'ACTIVE'
);

CREATE TABLE IF NOT EXISTS member_workout_progress (
  progress_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  follow_id UUID NOT NULL REFERENCES member_workout_follows(follow_id),
  completed_at TIMESTAMP NOT NULL DEFAULT now(),
  notes TEXT
);

-- =========================
-- RECOMMENDED INDEXES
-- =========================
CREATE INDEX IF NOT EXISTS idx_members_user_id ON members(user_id);
CREATE INDEX IF NOT EXISTS idx_members_qr_hash ON members(qr_static_hash);

CREATE INDEX IF NOT EXISTS idx_orders_status_created ON orders(status, created_at);
CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);

CREATE INDEX IF NOT EXISTS idx_member_memberships_member ON member_memberships(member_id, status, expire_at);
CREATE INDEX IF NOT EXISTS idx_session_usage_membership ON session_usage_logs(membership_id, used_at);

CREATE INDEX IF NOT EXISTS idx_otv_pass_token ON one_time_passes(token_hash);
CREATE INDEX IF NOT EXISTS idx_otv_pass_status ON one_time_passes(status, expires_at);

CREATE INDEX IF NOT EXISTS idx_attendance_scanned ON attendance_logs(scanned_at);
CREATE INDEX IF NOT EXISTS idx_attendance_user ON attendance_logs(user_id, scanned_at);
CREATE INDEX IF NOT EXISTS idx_attendance_pass ON attendance_logs(one_time_pass_id, scanned_at);

CREATE INDEX IF NOT EXISTS idx_access_state_gate ON user_access_state(gate_state);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_midtrans_order ON payments(midtrans_order_id);
CREATE INDEX IF NOT EXISTS idx_payments_midtrans_trx ON payments(midtrans_transaction_id);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_logs(created_at);

CREATE INDEX IF NOT EXISTS idx_workout_templates_creator ON workout_templates(created_by_user_id);
CREATE INDEX IF NOT EXISTS idx_workout_templates_coach ON workout_templates(coach_pt_id);
CREATE INDEX IF NOT EXISTS idx_workout_items_template ON workout_template_items(template_id, order_no);
CREATE INDEX IF NOT EXISTS idx_workout_follows_member ON member_workout_follows(member_id, followed_at);
CREATE INDEX IF NOT EXISTS idx_workout_progress_follow ON member_workout_progress(follow_id, completed_at);
