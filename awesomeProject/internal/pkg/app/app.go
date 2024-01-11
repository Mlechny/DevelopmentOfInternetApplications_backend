package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"awesomeProject/internal/app/config"
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/dsn"
	"awesomeProject/internal/app/repository"
)

type Application struct {
	repo   *repository.Repository
	config *config.Config
}

type GetLanguagesBack struct {
	Languages []ds.Language
	Name      string
}

func (a *Application) Run() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/html/*")

	r.GET("/languages", func(c *gin.Context) {
		name := c.Query("name")
		languages, err := a.repo.GetLanguageByName(name)

		if err != nil {
			log.Println("Can't get languages", err)
			c.Error(err)
			return
		}
		c.HTML(http.StatusOK, "all_codes.tmpl", GetLanguagesBack{
			Name:      name,
			Languages: languages,
		})

	})

	r.GET("/languages/:id", func(c *gin.Context) {
		id := c.Param("id")
		language, err := a.repo.GetLanguageByID(id)
		if err != nil {
			log.Printf("Can't get language by id %v", err)
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "code.tmpl", *language)
	})

	r.POST("/languages", func(c *gin.Context) {
		id := c.PostForm("delete")

		a.repo.DeleteLanguage(id)

		languages, err := a.repo.GetLanguageByName("")
		if err != nil {
			log.Println("Can't get languages", err)
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "all_codes.tmpl", GetLanguagesBack{
			Name:      "",
			Languages: languages,
		})
	})

	r.Static("/images", "./resources")
	r.Static("/styles", "./templates/css")
	err := r.Run("127.0.0.1:8080")
	if err != nil {
		return
	}
	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	return &app, nil
}
