package app

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"time"
	"web-service/internal/app/redis"
	"web-service/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"web-service/internal/app/config"
	"web-service/internal/app/dsn"
	"web-service/internal/app/repository"

	_ "web-service/docs"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
	redisClient *redis.Client
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	// Услуги - языки программирования
	api := r.Group("/api")
	{
		res := api.Group("/languages")
		{
			res.GET("/", app.WithAuthCheck(role.NotAuthorized, role.Student, role.Moderator), app.GetAllLanguages) // Список с поиском
			res.GET("/:id", app.WithAuthCheck(role.NotAuthorized, role.Student, role.Moderator), app.GetLanguage)  // Одна услуга
			res.DELETE("/:id", app.WithAuthCheck(role.Moderator), app.DeleteLanguage)                              // Удаление
			res.PUT("/:id", app.WithAuthCheck(role.Moderator), app.ChangeLanguage)                                 // Изменение
			res.POST("", app.WithAuthCheck(role.Moderator), app.AddLanguage)                                       // Добавление
			res.POST("/:id/add_to_form", app.WithAuthCheck(role.Student, role.Moderator), app.AddToForm)           // Добавление в заявку
		}

		// Заявки - формы
		n := api.Group("/forms")
		{
			n.GET("/", app.WithAuthCheck(role.Student, role.Moderator), app.GetAllForms)                          // Список (отфильтровать по дате формирования и статусу)
			n.GET("/:id", app.WithAuthCheck(role.Student, role.Moderator), app.GetForm)                           // Одна заявка
			n.PUT("", app.WithAuthCheck(role.Student, role.Moderator), app.UpdateForm)                            // Изменение (добавление комментариев)
			n.DELETE("", app.WithAuthCheck(role.Student, role.Moderator), app.DeleteForm)                         // Удаление
			n.DELETE("/delete_language/:id", app.WithAuthCheck(role.Student, role.Moderator), app.DeleteFromForm) // Изменеие (удаление услуг)
			n.PUT("/:id/change_github", app.WithAuthCheck(role.Student), app.ChangeGithub)                        //Изменение (добавление ссылки на гитхаб в м-м)
			n.PUT("/user_confirm", app.WithAuthCheck(role.Student, role.Moderator), app.UserConfirm)              // Сформировать создателем
			n.PUT("/:id/moderator_confirm", app.WithAuthCheck(role.Moderator), app.ModeratorConfirm)              // Завершить/отклонить модератором                                                                 //Добавление результата автотестирования (запрос к асинхронному сервису)
		}

		// Пользователи (авторизация и аутентификация)
		u := api.Group("/user")
		{
			u.POST("/sign_up", app.Register) //Авторизация
			u.POST("/login", app.Login)      //Регистрация
			u.GET("/logout", app.Logout)     //Выход из аккаунта
		}

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		err := r.Run(fmt.Sprintf("%s:%d", app.config.ServiceHost, app.config.ServicePort))
		if err != nil {
			return
		}

		log.Println("Server down")
	}
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

	app.minioClient, err = minio.New(app.config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	app.redisClient, err = redis.New(app.config.Redis)
	if err != nil {
		return nil, err
	}

	return &app, nil
}
