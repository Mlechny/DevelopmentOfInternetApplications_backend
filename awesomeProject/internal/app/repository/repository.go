package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"

	"awesomeProject/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetLanguageByID(id string) (*ds.Language, error) {
	language := &ds.Language{}
	err := r.db.Where("language_id = ?", id).First(language).Error
	if err != nil {
		return nil, err
	}

	return language, nil
}

func (r *Repository) GetAllLanguages() ([]ds.Language, error) {
	var languages []ds.Language

	err := r.db.Find(&languages).Error
	if err != nil {
		return nil, err
	}

	return languages, nil
}

func (r *Repository) GetLanguageByName(name string) ([]ds.Language, error) {
	var languages []ds.Language

	err := r.db.
		Where("LOWER(languages.name) LIKE ?", "%"+strings.ToLower(name)+"%").Where("is_deleted = ?", false).
		Find(&languages).Error

	if err != nil {
		return nil, err
	}

	return languages, nil
}

func (r *Repository) DeleteLanguage(id string) error {
	err := r.db.Exec("UPDATE languages SET is_deleted = ? WHERE language_id = ?", true, id).Error
	if err != nil {
		return err
	}

	return nil
}
