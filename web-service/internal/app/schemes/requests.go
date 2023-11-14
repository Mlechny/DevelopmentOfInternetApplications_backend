package schemes

import (
	"mime/multipart"
	"time"
	"web-service/internal/app/ds"
)

type LanguageRequest struct {
	LanguageId string `uri:"language_id" binding:"required,uuid"`
}

type GetAllLanguagesRequest struct {
	Name string `form:"name"`
}

type AddLanguageRequest struct {
	ds.Language
	Image *multipart.FileHeader `form:"image" json:"image" binding:"required"`
}

type ChangeLanguageRequest struct {
	LanguageId  string                `uri:"language_id" binding:"required,uuid"`
	Name        *string               `form:"name" json:"name" binding:"omitempty,max=50"`
	Subject     *string               `form:"subject" json:"subject" binding:"omitempty,max=70"`
	Task        *string               `form:"task" json:"task" binding:"omitempty,max=30"`
	Image       *multipart.FileHeader `form:"image" json:"image"`
	Description *string               `form:"description" json:"description" binding:"omitempty,max=1000"`
	Deadline    *time.Time            `form:"deadline" json:"deadline" time_format:"2006-01-02"`
}

type AddToFormRequest struct {
	LanguageId string `uri:"language_id" binding:"required,uuid"`
}

type GetAllFormsRequest struct {
	FormationDateStart *time.Time `form:"formation_date_start" json:"formation_date_start" time_format:"2006-01-02 15:04:05"`
	FormationDateEnd   *time.Time `form:"formation_date_end" json:"formation_date_end" time_format:"2006-01-02 15:04:05"`
	Status             string     `form:"status"`
}

type FormRequest struct {
	FormId string `uri:"form_id" binding:"required,uuid"`
}

type UpdateFormRequest struct {
	URI struct {
		FormId string `uri:"form_id" binding:"required,uuid"`
	}
	Comments string `form:"comments" json:"comments" binding:"required,max=300"`
}

type DeleteFromFormRequest struct {
	FormId     string `uri:"form_id" binding:"required,uuid"`
	LanguageId string `uri:"language_id" binding:"required,uuid"`
}

type UserConfirmRequest struct {
	FormId string `uri:"form_id" binding:"required,uuid"`
}

type ModeratorConfirmRequest struct {
	URI struct {
		FormId string `uri:"form_id" binding:"required,uuid"`
	}
	Status string `form:"status" json:"status" binding:"required"`
}
