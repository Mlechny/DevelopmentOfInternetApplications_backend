package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"web-service/internal/app/ds"
	"web-service/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.User{},
		&ds.Language{},
		&ds.Form{},
		&ds.Code{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
