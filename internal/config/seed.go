package config

import (
	"errors"
	"fiber-gorm-app/internal/models"
	"time"

	"gorm.io/gorm"
)

func SeedDevelopmentData(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	if err := seedTemplates(db); err != nil {
		return err
	}
	if err := seedTrainersAndSchedules(db); err != nil {
		return err
	}

	return nil
}

func seedTemplates(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Template{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	templates := []models.Template{
		{Name: "Beginner Full Body", Description: "3x/week full body plan untuk pemula."},
		{Name: "Push Pull Legs", Description: "Split PPL untuk intermediate (6 hari / minggu)."},
		{Name: "Fat Loss Starter", Description: "Kombinasi strength + cardio ringan."},
	}

	return db.Create(&templates).Error
}

func seedTrainersAndSchedules(db *gorm.DB) error {
	var trainerCount int64
	if err := db.Model(&models.Trainer{}).Count(&trainerCount).Error; err != nil {
		return err
	}
	if trainerCount > 0 {
		return nil
	}

	trainers := []models.Trainer{
		{Name: "Raka Putra", Bio: "Strength coach. Fokus compound lifts.", Specialty: "Strength", PhotoURL: ""},
		{Name: "Nadia Sari", Bio: "Fat loss & mobility. Friendly untuk pemula.", Specialty: "Fat Loss", PhotoURL: ""},
		{Name: "Dimas Pratama", Bio: "Hypertrophy programming.", Specialty: "Hypertrophy", PhotoURL: ""},
	}

	if err := db.Create(&trainers).Error; err != nil {
		return err
	}

	// Create schedules for each trainer
	var schedules []models.TrainerSchedule
	for _, t := range trainers {
		// seed Monday/Wednesday/Friday 18:00-20:00 as default
		schedules = append(schedules,
			models.TrainerSchedule{TrainerID: t.ID, DayOfWeek: 1, StartTime: "18:00", EndTime: "20:00", Location: "Revoks Gym"},
			models.TrainerSchedule{TrainerID: t.ID, DayOfWeek: 3, StartTime: "18:00", EndTime: "20:00", Location: "Revoks Gym"},
			models.TrainerSchedule{TrainerID: t.ID, DayOfWeek: 5, StartTime: "18:00", EndTime: "20:00", Location: "Revoks Gym"},
		)
	}

	// Ensure timestamps deterministic-ish
	now := time.Now()
	for i := range schedules {
		schedules[i].CreatedAt = now
		schedules[i].UpdatedAt = now
	}

	return db.Create(&schedules).Error
}
