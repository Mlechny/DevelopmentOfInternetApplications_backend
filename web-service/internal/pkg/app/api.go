package app

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"web-service/internal/app/ds"
	"web-service/internal/app/schemes"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func (app *Application) uploadImage(c *gin.Context, image *multipart.FileHeader, UUID string) (*string, error) {
	src, err := image.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	extension := filepath.Ext(image.Filename)
	if extension != ".jpg" && extension != ".jpeg" {
		return nil, fmt.Errorf("разрешены только jpeg изображения")
	}
	imageName := UUID + extension

	_, err = app.minioClient.PutObject(c, app.config.BucketName, imageName, src, image.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return nil, err
	}
	imageURL := fmt.Sprintf("%s/%s/%s", app.config.MinioEndpoint, app.config.BucketName, imageName)
	return &imageURL, nil
}

func (app *Application) getStudent() string {
	return "2d217868-ab6d-41fe-9b34-7809083a2e8a"
}

func (app *Application) getModerator() *string {
	moderatorId := "87d54d58-1e24-4cca-9c83-bd2523902729"
	return &moderatorId
}

func (app *Application) GetAllLanguages(c *gin.Context) {
	var request schemes.GetAllLanguagesRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	languages, err := app.repo.GetLanguageByName(request.Name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	draftForm, err := app.repo.GetDraftForm(app.getStudent())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := schemes.GetAllLanguagesResponse{DraftForm: nil, Languages: languages}
	if draftForm != nil {
		response.DraftForm = &schemes.FormShort{UUID: draftForm.UUID}
		containers, err := app.repo.GetCode(draftForm.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftForm.LanguageCount = len(containers)
	}
	c.JSON(http.StatusOK, response)
}

func (app *Application) GetLanguage(c *gin.Context) {
	var request schemes.LanguageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	language, err := app.repo.GetLanguageByID(request.LanguageId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if language == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("язык программирования не найден"))
		return
	}
	c.JSON(http.StatusOK, language)
}

func (app *Application) DeleteLanguage(c *gin.Context) {
	var request schemes.LanguageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	language, err := app.repo.GetLanguageByID(request.LanguageId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if language == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("язык программирования не найден"))
		return
	}
	language.IsDeleted = true
	if err := app.repo.SaveLanguage(language); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (app *Application) AddLanguage(c *gin.Context) {
	var request schemes.AddLanguageRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	language := request.Language
	if err := app.repo.AddLanguage(&language); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.Image != nil {
		imageURL, err := app.uploadImage(c, request.Image, language.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		language.ImageURL = imageURL
	}
	if err := app.repo.SaveLanguage(&language); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (app *Application) ChangeLanguage(c *gin.Context) {
	var request schemes.ChangeLanguageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	language, err := app.repo.GetLanguageByID(request.LanguageId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if language == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("язык программирования не найден"))
		return
	}

	if request.Name != nil {
		language.Name = *request.Name
	}
	if request.Image != nil {
		imageURL, err := app.uploadImage(c, request.Image, language.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		language.ImageURL = imageURL
	}
	if request.Subject != nil {
		language.Subject = *request.Subject
	}
	if request.Task != nil {
		language.Task = *request.Task
	}
	if request.Description != nil {
		language.Description = *request.Description
	}

	if err := app.repo.SaveLanguage(language); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, language)
}

func (app *Application) AddToForm(c *gin.Context) {
	var request schemes.AddToFormRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var err error

	language, err := app.repo.GetLanguageByID(request.LanguageId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if language == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("язык программирования не найден"))
		return
	}

	var form *ds.Form
	form, err = app.repo.GetDraftForm(app.getStudent())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		form, err = app.repo.CreateDraftForm(app.getStudent())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToForm(form.UUID, request.LanguageId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var languages []ds.Language
	languages, err = app.repo.GetCode(form.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllLanguagesResponse{Languages: languages})
}

func (app *Application) GetAllForms(c *gin.Context) {
	var request schemes.GetAllFormsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	forms, err := app.repo.GetAllForms(request.FormationDateStart, request.FormationDateEnd, request.Status)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	outputForms := make([]schemes.FormOutput, len(forms))
	for i, form := range forms {
		outputForms[i] = schemes.ConvertForm(&form)
	}
	c.JSON(http.StatusOK, schemes.AllFormsResponse{Forms: outputForms})
}

func (app *Application) GetForm(c *gin.Context) {
	var request schemes.FormRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	form, err := app.repo.GetFormById(request.FormId, app.getStudent())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}

	languages, err := app.repo.GetCode(request.FormId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.FormResponse{Form: schemes.ConvertForm(form), Languages: languages})
}

func (app *Application) UpdateForm(c *gin.Context) {
	var request schemes.UpdateFormRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	form, err := app.repo.GetFormById(request.URI.FormId, app.getStudent())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}
	form.Comments = request.Comments
	if app.repo.SaveForm(form); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.UpdateFormResponse{Form: schemes.ConvertForm(form)})
}

func (app *Application) DeleteForm(c *gin.Context) {
	var request schemes.FormRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	form, err := app.repo.GetFormById(request.FormId, app.getStudent())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}
	form.Status = ds.DELETED

	if err := app.repo.SaveForm(form); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) DeleteFromForm(c *gin.Context) {
	var request schemes.DeleteFromFormRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	form, err := app.repo.GetFormById(request.FormId, app.getStudent())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}
	if form.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя редактировать форму со статусом: %s", form.Status))
		return
	}

	if err := app.repo.DeleteFromForm(request.FormId, request.LanguageId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	languages, err := app.repo.GetCode(request.FormId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllLanguagesResponse{Languages: languages})
}

func (app *Application) UserConfirm(c *gin.Context) {
	var request schemes.UserConfirmRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	form, err := app.repo.GetFormById(request.FormId, app.getStudent())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}
	if form.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя сформировать форму со статусом %s", form.Status))
		return
	}
	form.Status = ds.FORMED
	now := time.Now()
	form.FormationDate = &now

	if err := app.repo.SaveForm(form); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) ModeratorConfirm(c *gin.Context) {
	var request schemes.ModeratorConfirmRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.Status != ds.COMPLETED && request.Status != ds.REJECTED {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("status %s not allowed", request.Status))
		return
	}

	form, err := app.repo.GetFormById(request.URI.FormId, app.getStudent())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}
	if form.Status != ds.FORMED {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", form.Status, request.Status))
		return
	}
	form.Status = request.Status
	form.ModeratorId = app.getModerator()
	if request.Status == ds.COMPLETED {
		now := time.Now()
		form.CompletionDate = &now
	}

	if err := app.repo.SaveForm(form); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
