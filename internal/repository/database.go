package repository

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"taskflow/internal/model"
)

func Connect(databaseURL string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к PostgreSQL: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Task{}); err != nil {
		log.Fatalf("Не удалось создать таблицы: %v", err)
	}
	return db

}
