package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"web-service/internal/app/ds"
	"web-service/internal/app/role"
	"web-service/internal/app/schemes"
)

// @Summary		Получить все формы
// @Tags		Формы
// @Description	Возвращает все формы с фильтрацией по статусу и дате формирования
// @Produce		json
// @Param		status query string false "статус формы"
// @Param		formation_date_start query string false "начальная дата формирования"
// @Param		formation_date_end query string false "конечная дата формирвания"
// @Success		200 {object} schemes.AllFormsResponse
// @Router		/api/forms [get]
func (app *Application) GetAllForms(c *gin.Context) {
	var request schemes.GetAllFormsRequest
	var err error
	if err = c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	var forms []ds.Form
	if userRole == role.Student {
		forms, err = app.repo.GetAllForms(&userId, request.FormationDateStart, request.FormationDateEnd, request.Status)
	} else {
		forms, err = app.repo.GetAllForms(nil, request.FormationDateStart, request.FormationDateEnd, request.Status)
	}
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

// @Summary		Получить одну форму
// @Tags		Формы
// @Description	Возвращает подробную информацию о форме и комментарий
// @Produce		json
// @Param		id path string true "id формы"
// @Success		200 {object} schemes.FormResponse
// @Router		/api/forms/{id} [get]
func (app *Application) GetForm(c *gin.Context) {
	var request schemes.FormRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	var form *ds.Form
	if userRole == role.Moderator {
		form, err = app.repo.GetFormById(request.FormId, nil)
	} else {
		form, err = app.repo.GetFormById(request.FormId, &userId)
	}
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

/*type SwaggerUpdateFormRequest struct {
	Comments string `json:"comments"`
}*/

// @Summary		Указать комментарий в форме
// @Tags		Формы
// @Description	Позволяет изменить комментарий в черновой форме и возвращает обновлённые данные
// @Access		json
// @Param		comments body schemes.UpdateFormRequest true "Комментарии"
// @Success		200
// @Router		/api/forms [put]
func (app *Application) UpdateForm(c *gin.Context) {
	var request schemes.UpdateFormRequest
	var err error
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}

	form.Comments = &request.Comments
	if app.repo.SaveForm(form); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Удалить черновую форму
// @Tags		Формы
// @Description	Удаляет черновую форму
// @Success		200
// @Router		/api/forms [delete]
func (app *Application) DeleteForm(c *gin.Context) {
	var err error
	var form *ds.Form
	userId := getUserId(c)
	form, err = app.repo.GetDraftForm(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}

	form.Status = ds.StatusDeleted

	if err := app.repo.SaveForm(form); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Удалить язык программирования из черновой формы
// @Tags		Формы
// @Description	Удалить язык программиования из черновой формы
// @Produce		json
// @Param		id path string true "id языка программирования"
// @Success		200
// @Router		/api/forms/delete_language/{id} [delete]
func (app *Application) DeleteFromForm(c *gin.Context) {
	var request schemes.DeleteFromFormRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}

	if err := app.repo.DeleteFromForm(form.UUID, request.LanguageId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Сформировать форму
// @Tags		Формы
// @Description	Сформировать или удалить форму пользователем
// @Success		200
// @Router		/api/forms/user_confirm [put]
func (app *Application) UserConfirm(c *gin.Context) {
	userId := getUserId(c)
	form, err := app.repo.GetDraftForm(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}
	if err := testingRequest(form.UUID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(`testing service is unavailable: {%s}`, err))
		return
	}

	testingStatus := ds.TestingStart
	form.Autotest = &testingStatus
	form.Status = ds.StatusFormed
	now := time.Now()
	form.FormationDate = &now

	if err := app.repo.SaveForm(form); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Подтвердить/отклонить форму
// @Tags		Формы
// @Description	Подтвердить или отклонить форму модератором
// @Param		id path string true "id формы"
// @Param		confirm body boolean true "подтвердить"
// @Success		200
// @Router		/api/forms/{id}/moderator_confirm [put]
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

	userId := getUserId(c)
	form, err := app.repo.GetFormById(request.URI.FormId, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("заявка не найдена"))
		return
	}
	if form.Status != ds.StatusFormed {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", form.Status, ds.StatusFormed))
		return
	}

	if *request.Confirm {
		form.Status = ds.StatusCompleted
	} else {
		form.Status = ds.StatusRejected
	}
	now := time.Now()
	form.CompletionDate = &now
	moderator, err := app.repo.GetUserById(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	form.Moderator = moderator

	if err := app.repo.SaveForm(form); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) Testing(c *gin.Context) {
	var request schemes.TestingReq
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.Token != app.config.Token {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	form, err := app.repo.GetFormById(request.URI.FormId, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}

	var testingStatus string
	if *request.TestingStatus {
		testingStatus = ds.TestingSuccess
	} else {
		testingStatus = ds.TestingFailure
	}
	form.Autotest = &testingStatus

	if err := app.repo.SaveForm(form); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Указать ссылку на гитхаб
// @Tags		Коды
// @Description	Позволяет изменить ссылку на гитхаб в таблице м-м и возвращает обновленные данные
// @Access      json
// @Produce     json
// @Param		id path string true "id формы"
// @Param		github body schemes.ChangeGithubRequest true "Гитхаб"
// @Success		200 {object} schemes.CodeResponse
// @Router		/api/forms/{id}/change_github [put]
func (app *Application) ChangeGithub(c *gin.Context) {
	var request schemes.ChangeGithubRequest

	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	fmt.Println(userId, userRole)

	code, err := app.repo.GetCodeByFormId(request.FormId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	form, err := app.repo.GetFormById(code.FormId, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if form == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("форма не найдена"))
		return
	}

	if form.StudentId != userId {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("изменить поле может только создатель формы"))
		return
	}

	code.Github = &request.Github

	if err := app.repo.SaveCode(code); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.CodeResponse{Code: schemes.ConvertCode(code)})
}
