package repository

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"web-service/internal/app/ds"
)

func (r *Repository) GetAllForms(formationDateStart, formationDateEnd *time.Time, status string) ([]ds.Form, error) {
	var forms []ds.Form
	query := r.db.Preload("Student").Preload("Moderator").
		Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
		Where("status != ?", ds.DELETED)

	if formationDateStart != nil && formationDateEnd != nil {
		query = query.Where("formation_date BETWEEN ? AND ?", *formationDateStart, *formationDateEnd)
	} else if formationDateStart != nil {
		query = query.Where("formation_date >= ?", *formationDateStart)
	} else if formationDateEnd != nil {
		query = query.Where("formation_date <= ?", *formationDateEnd)
	}
	if err := query.Find(&forms).Error; err != nil {
		return nil, err
	}
	return forms, nil
}

func (r *Repository) GetDraftForm(studentId string) (*ds.Form, error) {
	form := &ds.Form{}
	err := r.db.First(form, ds.Form{Status: ds.DRAFT, StudentId: studentId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return form, nil
}

func (r *Repository) CreateDraftForm(studentId string) (*ds.Form, error) {
	form := &ds.Form{CreationDate: time.Now(), StudentId: studentId, Status: ds.DRAFT}
	err := r.db.Create(form).Error
	if err != nil {
		return nil, err
	}
	return form, nil
}

func (r *Repository) GetFormById(formId, studentId string) (*ds.Form, error) {
	form := &ds.Form{}
	err := r.db.Preload("Moderator").Preload("Student").
		Where("status != ?", ds.DELETED).
		First(form, ds.Form{UUID: formId, StudentId: studentId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return form, nil
}

func (r *Repository) GetCode(formId string) ([]ds.Language, error) {
	var languages []ds.Language

	err := r.db.Table("codes").
		Select("languages.*").
		Joins("JOIN languages ON codes.language_id = languages.uuid").
		Where(ds.Code{FormId: formId}).
		Scan(&languages).Error

	if err != nil {
		return nil, err
	}
	return languages, nil
}

func (r *Repository) SaveForm(form *ds.Form) error {
	err := r.db.Save(form).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteFromForm(formId, languageId string) error {
	err := r.db.Delete(&ds.Code{FormId: formId, LanguageId: languageId}).Error
	if err != nil {
		return err
	}
	return nil
}
