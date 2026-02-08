package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
	)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		ClauseBuilders: map[string]clause.ClauseBuilder{
			"LIMIT": func(c clause.Clause, builder clause.Builder) {
				lim, ok := c.Expression.(clause.Limit)
				if !ok {
					return
				}
				wrote := false
				if lim.Limit != nil {
					builder.WriteString("LIMIT ")
					builder.WriteString(strconv.Itoa(*lim.Limit))
					wrote = true
				}
				if lim.Offset > 0 {
					if wrote {
						builder.WriteByte(' ')
					}
					builder.WriteString("OFFSET ")
					builder.WriteString(strconv.Itoa(lim.Offset))
				}
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to init gorm: %v", err)
	}

	log.Println("✅ Database connected successfully")

	// Run raw SQL migration instead of AutoMigrate
	log.Println("📦 Running SQL migration...")
	if err := runMigration(sqlDB); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ All tables migrated successfully!")
}

func runMigration(db *sql.DB) error {
	migration := `
	-- Extensions
	CREATE EXTENSION IF NOT EXISTS "pgcrypto";
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	-- =====================
	-- Phase 1: Base tables
	-- =====================
	CREATE TABLE IF NOT EXISTS users (
		user_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		role          VARCHAR(20) NOT NULL,
		email         VARCHAR(255),
		phone         VARCHAR(50),
		password_hash VARCHAR(255) NOT NULL,
		full_name     VARCHAR(255) NOT NULL,
		photo_url     VARCHAR(500),
		status        VARCHAR(20) NOT NULL DEFAULT 'active',
		created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE email IS NOT NULL;
	CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone ON users(phone) WHERE phone IS NOT NULL;

	CREATE TABLE IF NOT EXISTS scanner_devices (
		scanner_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name         VARCHAR(100) NOT NULL,
		location     VARCHAR(200),
		api_key_hash VARCHAR(255) NOT NULL,
		is_active    BOOLEAN NOT NULL DEFAULT TRUE,
		created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS exercise_catalog (
		exercise_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name           VARCHAR(200) NOT NULL,
		primary_muscle VARCHAR(100),
		equipment      VARCHAR(100),
		video_url      VARCHAR(500),
		image_url      VARCHAR(500),
		is_active      BOOLEAN NOT NULL DEFAULT TRUE,
		created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS membership_plans (
		plan_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name          VARCHAR(100) NOT NULL,
		price_amount  DECIMAL NOT NULL,
		duration_days INT,
		session_quota INT,
		is_active     BOOLEAN NOT NULL DEFAULT TRUE,
		created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS one_time_products (
		otv_product_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name              VARCHAR(100) NOT NULL,
		price_amount      DECIMAL NOT NULL,
		validity_minutes  INT NOT NULL,
		is_active         BOOLEAN NOT NULL DEFAULT TRUE,
		created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	-- ============================
	-- Phase 2: User-dependent
	-- ============================
	CREATE TABLE IF NOT EXISTS members (
		member_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id          UUID NOT NULL UNIQUE REFERENCES users(user_id),
		legacy_member_no BIGINT,
		legacy_record_id BIGINT,
		gender           VARCHAR(10),
		birth_date       DATE,
		address          TEXT,
		phone_override   VARCHAR(50),
		qr_static_hash   VARCHAR(255) NOT NULL UNIQUE,
		created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS personal_trainers (
		pt_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id          UUID NOT NULL UNIQUE REFERENCES users(user_id),
		bio              TEXT,
		experience_years INT,
		specialties      TEXT,
		base_rate_amount DECIMAL,
		rating_avg       DECIMAL,
		rating_count     INT,
		is_active        BOOLEAN NOT NULL DEFAULT TRUE,
		created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS user_access_state (
		user_id         UUID PRIMARY KEY REFERENCES users(user_id),
		gate_state      VARCHAR(20) NOT NULL DEFAULT 'OUTSIDE',
		last_scanned_at TIMESTAMP,
		last_scanner_id UUID REFERENCES scanner_devices(scanner_id)
	);

	CREATE TABLE IF NOT EXISTS notifications (
		notification_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id         UUID NOT NULL REFERENCES users(user_id),
		type            VARCHAR(50) NOT NULL,
		title           VARCHAR(255) NOT NULL,
		body            TEXT,
		payload_json    JSONB,
		read_at         TIMESTAMP,
		created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		audit_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		actor_user_id UUID REFERENCES users(user_id),
		action        VARCHAR(100) NOT NULL,
		entity_type   VARCHAR(100),
		entity_id     UUID,
		payload_json  JSONB,
		created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS workout_templates (
		template_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title              VARCHAR(200) NOT NULL,
		description        TEXT,
		focus_muscle       VARCHAR(100),
		level              VARCHAR(50),
		created_by_user_id UUID REFERENCES users(user_id),
		coach_pt_id        UUID REFERENCES personal_trainers(pt_id),
		is_active          BOOLEAN NOT NULL DEFAULT TRUE,
		created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	-- =====================
	-- Phase 3: Orders
	-- =====================
	CREATE TABLE IF NOT EXISTS orders (
		order_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		channel      VARCHAR(50) NOT NULL,
		buyer_user_id UUID REFERENCES users(user_id),
		status       VARCHAR(20) NOT NULL DEFAULT 'PENDING',
		total_amount DECIMAL NOT NULL,
		currency     VARCHAR(10) NOT NULL DEFAULT 'IDR',
		created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		expired_at   TIMESTAMP
	);

	-- =====================
	-- Phase 4: Order Items
	-- =====================
	CREATE TABLE IF NOT EXISTS order_items (
		order_item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_id      UUID NOT NULL REFERENCES orders(order_id),
		item_type     VARCHAR(50) NOT NULL,
		ref_id        UUID NOT NULL,
		qty           INT NOT NULL DEFAULT 1,
		unit_price    DECIMAL NOT NULL,
		subtotal      DECIMAL NOT NULL
	);

	-- ============================
	-- Phase 5: Order-dependent
	-- ============================
	CREATE TABLE IF NOT EXISTS member_memberships (
		membership_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		member_id              UUID NOT NULL REFERENCES members(member_id),
		plan_id                UUID NOT NULL REFERENCES membership_plans(plan_id),
		start_at               TIMESTAMP,
		expire_at              TIMESTAMP,
		status                 VARCHAR(20) NOT NULL,
		total_sessions         INT,
		used_sessions          INT,
		remaining_sessions     INT,
		activated_by_order_id  UUID REFERENCES orders(order_id),
		created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS one_time_passes (
		one_time_pass_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_id         UUID NOT NULL REFERENCES orders(order_id),
		otv_product_id   UUID NOT NULL REFERENCES one_time_products(otv_product_id),
		token_hash       VARCHAR(255) NOT NULL UNIQUE,
		status           VARCHAR(20) NOT NULL DEFAULT 'UNUSED',
		issued_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		expires_at       TIMESTAMP,
		used_at          TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS payments (
		payment_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_id                UUID NOT NULL REFERENCES orders(order_id),
		provider                VARCHAR(50) NOT NULL DEFAULT 'MIDTRANS',
		midtrans_order_id       VARCHAR(100) NOT NULL UNIQUE,
		midtrans_transaction_id VARCHAR(100) UNIQUE,
		payment_type            VARCHAR(50),
		transaction_status      VARCHAR(50),
		fraud_status            VARCHAR(50),
		gross_amount            DECIMAL NOT NULL,
		currency                VARCHAR(10) NOT NULL DEFAULT 'IDR',
		actions_json            JSONB,
		va_numbers_json         JSONB,
		metadata_json           JSONB,
		status                  VARCHAR(20) NOT NULL DEFAULT 'PENDING',
		paid_at                 TIMESTAMP,
		expires_at              TIMESTAMP,
		raw_response_json       JSONB,
		created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS payment_webhook_events (
		event_id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		provider                VARCHAR(50) NOT NULL DEFAULT 'MIDTRANS',
		provider_event_id       VARCHAR(100) UNIQUE,
		order_id                UUID REFERENCES orders(order_id),
		midtrans_order_id       VARCHAR(100),
		midtrans_transaction_id VARCHAR(100),
		raw_payload_json        JSONB NOT NULL,
		received_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		processed_at            TIMESTAMP,
		process_status          VARCHAR(20) NOT NULL DEFAULT 'RECEIVED'
	);

	-- ============================
	-- Phase 6: PT-related
	-- ============================
	CREATE TABLE IF NOT EXISTS pt_reviews (
		review_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		pt_id      UUID NOT NULL REFERENCES personal_trainers(pt_id),
		member_id  UUID NOT NULL REFERENCES members(member_id),
		rating     INT,
		comment    TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS pt_availability_rules (
		rule_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		pt_id      UUID NOT NULL REFERENCES personal_trainers(pt_id),
		weekday    VARCHAR(10) NOT NULL,
		start_time TIME NOT NULL,
		end_time   TIME NOT NULL,
		is_active  BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	-- ============================
	-- Phase 7: Workout
	-- ============================
	CREATE TABLE IF NOT EXISTS workout_template_items (
		item_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		template_id  UUID NOT NULL REFERENCES workout_templates(template_id),
		exercise_id  UUID NOT NULL REFERENCES exercise_catalog(exercise_id),
		order_no     INT,
		sets         INT,
		reps         VARCHAR(50),
		rest_seconds INT,
		tempo        VARCHAR(20),
		rpe          INT,
		notes        TEXT
	);

	CREATE TABLE IF NOT EXISTS member_workout_follows (
		follow_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		member_id   UUID NOT NULL REFERENCES members(member_id),
		template_id UUID NOT NULL REFERENCES workout_templates(template_id),
		followed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		status      VARCHAR(20) NOT NULL DEFAULT 'ACTIVE'
	);

	CREATE TABLE IF NOT EXISTS member_workout_progress (
		progress_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		follow_id    UUID NOT NULL REFERENCES member_workout_follows(follow_id),
		completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		notes        TEXT
	);

	-- ============================
	-- Phase 8: Logs
	-- ============================
	CREATE TABLE IF NOT EXISTS session_usage_logs (
		usage_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		membership_id UUID NOT NULL REFERENCES member_memberships(membership_id),
		used_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		used_by       VARCHAR(100),
		notes         VARCHAR(500)
	);

	CREATE TABLE IF NOT EXISTS attendance_logs (
		attendance_log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		scanned_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		scanner_id        UUID NOT NULL REFERENCES scanner_devices(scanner_id),
		user_id           UUID REFERENCES users(user_id),
		one_time_pass_id  UUID REFERENCES one_time_passes(one_time_pass_id),
		scanned_token_hash VARCHAR(255),
		direction         VARCHAR(10),
		result            VARCHAR(20) NOT NULL,
		reject_reason     VARCHAR(255)
	);
	`

	_, err := db.Exec(migration)
	return err
}

func GetDB() *gorm.DB {
	return DB
}
