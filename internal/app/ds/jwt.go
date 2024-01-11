package ds

import (
	"github.com/golang-jwt/jwt"
	"web-service/internal/app/role"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserUUID string    `json:"user_uuid"`
	Role     role.Role `json:"role"`
	Login    string    `json:"login"`
}
