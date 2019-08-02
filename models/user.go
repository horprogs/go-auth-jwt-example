package models

import (
	"auth-go-example/lib/customErrors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_JWT_KEY = []byte("ANB4xetePU2hHFM1tUF5FTsQbqdxDY6A")

type User struct {
	gorm.Model
	UserID   string `gorm:"primary_key:true"json:"id"`
	Name     string `json:"name"`
	Password string
	Email    string
	Phone    string
}

type AuthTable struct {
	gorm.Model
	UserID       string
	RefreshToken string
}

type AuthObject struct {
	Password string
	Login    string
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJWT(userId string) string {
	var mySigningKey = SECRET_JWT_KEY

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
		"user_id": userId,
	}

	tokenString, _ := token.SignedString(mySigningKey)

	return tokenString
}

func generateRefreshToken() string {
	return randstr.Hex(16)
}

func GetToken(login string, password string) (map[string]string, error) {
	var foundUser User
	res := map[string]string{
		"accessToken":  "",
		"refreshToken": "",
	}

	db := Database()
	defer db.Close()

	db.Where(&User{Email: login}).First(&foundUser)

	if foundUser.ID == 0 {
		return res, &customErrors.RespError{Code: 404, Message: "User not found"}
	}

	isValid := CheckPasswordHash(password, foundUser.Password)

	if isValid != true {
		return res, &customErrors.RespError{Code: 403, Message: "Password wrong"}
	}

	accessToken := generateJWT(foundUser.UserID)
	refreshToken := generateRefreshToken()

	authTable := AuthTable{
		UserID: foundUser.UserID,
	}

	db.First(&authTable)

	authTable.RefreshToken = refreshToken

	db.Save(&authTable)

	res["accessToken"] = accessToken
	res["refreshToken"] = refreshToken

	return res, nil
}

func RefreshToken(userId string, refreshToken string) (map[string]string, error) {
	var authUser AuthTable
	res := map[string]string{
		"accessToken":  "",
		"refreshToken": "",
	}

	db := Database()
	defer db.Close()

	db.Where(&AuthTable{UserID: userId}).First(&authUser)

	if refreshToken != authUser.RefreshToken {
		return res, &customErrors.RespError{Code: 500, Message: "RefreshToken invalid"}
	}

	newRefreshToken := randstr.Hex(16)
	newAccessToken := generateJWT(userId)

	res["accessToken"] = newAccessToken
	res["refreshToken"] = newRefreshToken

	authUser.RefreshToken = newRefreshToken

	db.Save(&authUser)

	return res, nil
}
