package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"web-service/internal/app/ds"
	"web-service/internal/app/schemes"
)

// @Summary		Получить все языки программирования
// @Tags		Языки программирования
// @Description	Возвращает все доступные языки программирования с опциональной фильтрацией по названию
// @Produce		json
// @Param		name query string false "Название для фильтрации"
// @Success		200 {object} schemes.GetAllLanguagesResponse
// @Router		/api/languages [get]
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
	response := schemes.GetAllLanguagesResponse{DraftForm: nil, Languages: languages}
	if userId, exists := c.Get("userId"); exists {
		draftForm, err := app.repo.GetDraftForm(userId.(string))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if draftForm != nil {
			response.DraftForm = &draftForm.UUID
		}
	}
	c.JSON(http.StatusOK, response)
}

// @Summary		Получить один язык программирования
// @Tags		Языки программирования
// @Description	Возвращает более подробную информацию об одном языке программирования
// @Produce		json
// @Param		id path string true "id языка программирования"
// @Success		200 {object} ds.Language
// @Router		/api/languages/{id} [get]
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

// @Summary		Удалить язык программирования
// @Tags		Языки программирования
// @Description	Удаляет язык программирования по id
// @Param		id path string true "id языка программирования"
// @Success		200
// @Router		/api/languages/{id} [delete]
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

	if language.ImageURL != nil {
		if err := app.deleteImage(c, language.UUID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	language.ImageURL = nil
	language.IsDeleted = true
	if err := app.repo.SaveLanguage(language); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Добавить язык программирования
// @Tags		Языки программирования
// @Description	Добавить новый язык программирования
// @Accept		mpfd
// @Param     	image formData file false "Изображение"
// @Param     	name formData string true "Название" format:"string" maxLength:20
// @Param     	subject formData string true "Предмет" format:"string" maxLength:70
// @Param     	task formData int true "Задание" format:"string" maxLength:30
// @Param     	description formData string true "Описание задания" format:"string" maxLength:1000
// @Success		200
// @Router		/api/languages [post]
func (app *Application) AddLanguage(c *gin.Context) {
	var request schemes.AddLanguageRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	language := request.Language ///?
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

	c.JSON(http.StatusCreated, language.UUID)
}

// @Summary		Изменить язык программирования
// @Tags		Языки программирования
// @Description	Изменить данные полей об языке программирования
// @Accept		mpfd
// @Produce		json
// @Param		id path string true "Идентификатор языка программирования" format:"uuid"
// @Param		name formData string false "Название" format:"string" maxLength:20
// @Param		subject formData string false "Предмет" format:"string" maxLength:70
// @Param		task formData int false "Задание" format:"string" maxLength:30
// @Param		image formData file false "Изображение"
// @Param		description formData string false "Описание задания" format:"string" maxLength:1000
// @Router		/api/languages/{id} [put]
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
		if language.ImageURL != nil {
			if err := app.deleteImage(c, language.UUID); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
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

	c.Status(http.StatusOK)
}

// @Summary		Добавить в форму
// @Tags		Языки программирования
// @Description	Добавить выбранный язык программирования в черновик формы
// @Produce		json
// @Param		id path string true "id языка программирования"
// @Success		201 {object} schemes.AddToFormResp
// @Router		/api/languages/{id}/add_to_form [post]
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
	userId := getUserId(c)
	form, err = app.repo.GetDraftForm(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		form, err = app.repo.CreateDraftForm(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToForm(form.UUID, request.LanguageId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
