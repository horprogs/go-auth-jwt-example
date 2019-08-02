package middlewares

import (
	"auth-go-example/models"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
)

type MyJWTClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

func getPayload(tokenString string) *MyJWTClaims {
	token, _ := jwt.ParseWithClaims(tokenString, &MyJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return models.SECRET_JWT_KEY, nil
	})

	claims := token.Claims.(*MyJWTClaims)

	return claims
}

func isTokenValid(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return models.SECRET_JWT_KEY, nil
	})

	if token.Valid {
		return true
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return false
		} else {
			return false
		}
	} else {
		return false
	}
}

func AuthMiddleware(ctx *context.Context) {
	if ctx.Request.Header["Authorization"] != nil {
		token := ctx.Request.Header["Authorization"][0]

		userId := getPayload(token).UserId


		ctx.Output.Session("userIdForRefresh", userId)

		if isTokenValid(token) {
			ctx.Input.SetParam("userId", userId)
		}
	}
}
