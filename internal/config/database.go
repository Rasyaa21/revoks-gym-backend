package config

import (
	"fiber-gorm-app/internal/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

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

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Auto migrate
	err = DB.AutoMigrate(
		&models.User{},
		&models.Membership{},
		&models.AttendanceLog{},
		&models.WorkoutProgress{},
		&models.Template{},
		&models.UserTemplateFollow{},
		&models.Target{},
		&models.TargetProgress{},
		&models.Trainer{},
		&models.TrainerSchedule{},
		&models.Notification{},
		&models.UserSetting{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migrated successfully")

	if shouldSeed() {
		if err := SeedDevelopmentData(DB); err != nil {
			log.Printf("Seed skipped/failed: %v", err)
		} else {
			log.Println("Database seeded successfully")
		}
	}
}

func shouldSeed() bool {
	return os.Getenv("DB_SEED") == "true"
}

func GetDB() *gorm.DB {
	return DB
}
