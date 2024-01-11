package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"web-service/internal/app/config"
	"web-service/internal/app/dsn"
	"web-service/internal/app/repository"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	r.GET("/api/languages", app.GetAllLanguages)
	r.GET("/api/languages/:language_id", app.GetLanguage)
	r.DELETE("/api/languages/:language_id", app.DeleteLanguage)
	r.PUT("/api/languages/:language_id", app.ChangeLanguage)
	r.POST("/api/languages", app.AddLanguage)
	r.POST("/api/languages/:language_id/add_to_form", app.AddToForm)

	r.GET("/api/forms", app.GetAllForms)
	r.GET("/api/forms/:form_id", app.GetForm)
	r.PUT("/api/forms/:form_id/update", app.UpdateForm)
	r.DELETE("/api/forms/:form_id", app.DeleteForm)
	r.DELETE("/api/forms/:form_id/delete_language/:language_id", app.DeleteFromForm)
	r.PUT("/api/forms/:form_id/user_confirm", app.UserConfirm)
	r.PUT("/api/forms/:form_id/moderator_confirm", app.ModeratorConfirm)

	r.Static("/image", "./resources")
	r.Static("/css", "./templates/css")
	err := r.Run("127.0.0.1:8080")
	if err != nil {
		return
	}
	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	loc, _ := time.LoadLocation("Europe/Moscow")
	time.Local = loc
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	app.minioClient, err = minio.New(app.config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			log.Println(err.Err)
		}
		lastError := c.Errors.Last()
		if lastError != nil {
			switch c.Writer.Status() {
			case http.StatusBadRequest:
				c.JSON(-1, gin.H{"error": "wrong request"})
			case http.StatusNotFound:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			case http.StatusMethodNotAllowed:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			default:
				c.Status(-1)
			}
		}
	}
}
