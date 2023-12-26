package schemes

import (
	"time"
	"web-service/internal/app/ds"
)

type AllLanguagesResponse struct {
	Languages []ds.Language `json:"languages"`
}

type FormShort struct {
	UUID          string `json:"uuid"`
	LanguageCount int    `json:"language_count"`
}

type GetAllLanguagesResponse struct {
	DraftForm *FormShort    `json:"draft_form"`
	Languages []ds.Language `json:"languages"`
}

type AllFormsResponse struct {
	Forms []FormOutput `json:"forms"`
}

type FormResponse struct {
	Form      FormOutput    `json:"form"`
	Languages []ds.Language `json:"languages"`
}

type UpdateFormResponse struct {
	Form FormOutput `json:"forms"`
}

type FormOutput struct {
	UUID           string  `json:"uuid"`
	Status         string  `json:"status"`
	CreationDate   string  `json:"creation_date"`
	FormationDate  *string `json:"formation_date"`
	CompletionDate *string `json:"completion_date"`
	Moderator      *string `json:"moderator"`
	Student        string  `json:"student"`
	Comments       *string `json:"comments"`
	SendingStatus  *string `json:"sending_status"`
}

func ConvertForm(form *ds.Form) FormOutput {
	output := FormOutput{
		UUID:         form.UUID,
		Status:       form.Status,
		CreationDate: form.CreationDate.Format("2006-01-02 15:04:05"),
		Comments:     form.Comments,
		Student:      form.Student.Login,
	}

	if form.FormationDate != nil {
		formationDate := form.FormationDate.Format("2006-01-02 15:04:05")
		output.FormationDate = &formationDate
	}

	if form.CompletionDate != nil {
		completionDate := form.CompletionDate.Format("2006-01-02 15:04:05")
		output.CompletionDate = &completionDate
	}

	if form.Moderator != nil {
		output.Moderator = &form.Moderator.Login
	}

	return output
}

type AddToFormResp struct {
	LanguagesCount int64 `json:"language_count"`
}

type LoginResp struct {
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
}

type SwaggerLoginResp struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type RegisterResp struct {
	Ok bool `json:"ok"`
}
