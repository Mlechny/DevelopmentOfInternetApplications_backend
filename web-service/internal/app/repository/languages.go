package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"web-service/internal/app/ds"
)

func (r *Repository) GetLanguageByID(id string) (*ds.Language, error) {
	language := &ds.Language{UUID: id}
	err := r.db.First(language, "is_deleted = ?", false).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return language, nil
}

func (r *Repository) AddLanguage(language *ds.Language) error {
	err := r.db.Create(&language).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetLanguageByName(Name string) ([]ds.Language, error) {
	var languages []ds.Language

	err := r.db.
		Where("LOWER(name) LIKE ?", "%"+strings.ToLower(Name)+"%").Where("is_deleted = ?", false).
		Find(&languages).Error

	if err != nil {
		return nil, err
	}

	return languages, nil
}

func (r *Repository) SaveLanguage(language *ds.Language) error {
	err := r.db.Save(language).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddToForm(formId, languageId string) error {
	code := ds.Code{FormId: formId, LanguageId: languageId}
	err := r.db.Create(&code).Error
	if err != nil {
		return err
	}
	return nil
}
