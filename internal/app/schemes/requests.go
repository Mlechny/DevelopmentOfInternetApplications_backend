package schemes

import (
	"mime/multipart"
	"time"
	"web-service/internal/app/ds"
)

type LanguageRequest struct {
	LanguageId string `uri:"id" binding:"required,uuid"`
}

type GetAllLanguagesRequest struct {
	Name string `form:"name"`
}

type AddLanguageRequest struct {
	ds.Language
	Image *multipart.FileHeader `form:"image" json:"image"`
}

type ChangeLanguageRequest struct {
	LanguageId  string                `uri:"id" binding:"required,uuid"`
	Name        *string               `form:"name" json:"name" binding:"omitempty,max=50"`
	Subject     *string               `form:"subject" json:"subject" binding:"omitempty,max=70"`
	Task        *string               `form:"task" json:"task" binding:"omitempty,max=30"`
	Image       *multipart.FileHeader `form:"image" json:"image"`
	Description *string               `form:"description" json:"description" binding:"omitempty,max=1000"`
}

type ChangeGithubRequest struct {
	FormId string `uri:"id" binding:"required,uuid"`
	Github string `form:"github" json:"github" binding:"required,max=50"`
}

type AddToFormRequest struct {
	LanguageId string `uri:"id" binding:"required,uuid"`
}

type GetAllFormsRequest struct {
	FormationDateStart *time.Time `form:"formation_date_start" json:"formation_date_start" time_format:"2006-01-02 15:04"`
	FormationDateEnd   *time.Time `form:"formation_date_end" json:"formation_date_end" time_format:"2006-01-02 15:04"`
	Status             string     `form:"status"`
}

type FormRequest struct {
	FormId string `uri:"id" binding:"required,uuid"`
}

type UpdateFormRequest struct {
	Comments string `form:"comments" json:"comments" binding:"required,max=300"`
}

type DeleteFromFormRequest struct {
	LanguageId string `uri:"id" binding:"required,uuid"`
}

type ModeratorConfirmRequest struct {
	URI struct {
		FormId string `uri:"id" binding:"required,uuid"`
	}
	Confirm *bool `form:"confirm" binding:"required"`
}

type LoginReq struct {
	Login    string `form:"login" binding:"required,max=30"`
	Password string `form:"password" binding:"required,max=30"`
}

type RegisterReq struct {
	Login    string `form:"login" binding:"required,max=30"`
	Password string `form:"password" binding:"required,max=30"`
}