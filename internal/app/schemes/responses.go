package schemes

import (
	"web-service/internal/app/ds"
)

type AllLanguagesResponse struct {
	Languages []ds.Language `json:"languages"`
}

type GetAllLanguagesResponse struct {
	DraftForm *string       `json:"draft_form"`
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

type CodeResponse struct {
	Code CodeOutput `json:"code"`
}

type CodeOutput struct {
	LanguageId string  `json:"language_id"`
	FormId     string  `json:"form_id"`
	Github     *string `json:"github"`
}

func ConvertCode(code *ds.Code) CodeOutput {
	output := CodeOutput{
		LanguageId: code.LanguageId,
		FormId:     code.FormId,
		Github:     code.Github,
	}
	return output
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
	Autotest       *string `json:"autotest"`
}

func ConvertForm(form *ds.Form) FormOutput {
	output := FormOutput{
		UUID:         form.UUID,
		Status:       form.Status,
		CreationDate: form.CreationDate.Format("2006-01-02 15:04:05"),
		Comments:     form.Comments,
		Student:      form.Student.Login,
		Autotest:     form.Autotest,
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

type AuthResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
