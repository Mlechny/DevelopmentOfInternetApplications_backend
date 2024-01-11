package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"

	"net/http"
	"strings"
	"time"
	"web-service/internal/app/ds"
	"web-service/internal/app/role"
	"web-service/internal/app/schemes"
)

// @Summary		Регистрация
// @Tags		Авторизация и аутентификация
// @Description	Регистрация нового пользователя
// @Accept		json
// @Param		user_credentials body schemes.RegisterReq true "login and password"
// @Success		200
// @Router		/api/user/sign_up [post]
func (app *Application) Register(c *gin.Context) {
	request := &schemes.RegisterReq{}
	if err := c.ShouldBind(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	existing_user, err := app.repo.GetUserByLogin(request.Login)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if existing_user != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := ds.User{
		Role:     role.Student,
		Login:    request.Login,
		Password: generateHashString(request.Password),
	}
	if err := app.repo.AddUser(&user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Авторизация
// @Tags		Авторизация и аутентификация
// @Description	Авторизует пользователя по логиню, паролю и отдаёт jwt токен для дальнейших запросов
// @Accept		json
// @Produce		json
// @Param		user_credentials body schemes.LoginReq true "login and password"
// @Success		200 {object} schemes.AuthResp
// @Router		/api/user/login [post]
// @Consumes     json
func (app *Application) Login(c *gin.Context) {
	JWTConfig := app.config.JWT
	request := &schemes.LoginReq{}
	if err := c.ShouldBind(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := app.repo.GetUserByLogin(request.Login)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if user.Password != generateHashString(request.Password) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	token := jwt.NewWithClaims(JWTConfig.SigningMethod, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(JWTConfig.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserUUID: user.UUID,
		Role:     user.Role,
		Login:    user.Login,
	})
	if token == nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
		return
	}

	strToken, err := token.SignedString([]byte(JWTConfig.Token))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
		return
	}

	c.JSON(http.StatusOK, schemes.AuthResp{
		AccessToken: strToken,
		TokenType:   "Bearer",
	})
}

// @Summary		Выход из аккаунта
// @Tags		Авторизация и аутентификация
// @Description	Осуществляет выход из аккаунта
// @Accept		json
// @Success		200
// @Router		/api/user/loguot [get]
func (app *Application) Logout(c *gin.Context) {
	jwtStr := c.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf(jwtStr))
		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.config.JWT.Token), nil
	})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	err = app.redisClient.WriteJWTToBlacklist(c.Request.Context(), jwtStr, app.config.JWT.ExpiresIn)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
