package controllers

import (
	"auth-go-example/lib/customErrors"
	"auth-go-example/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

type ResponseUser struct {
	UserID  string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	IsAdmin bool   `json:"isAdmin"`
}

// @router /get-token [post]
func (u *UserController) GetToken() {
	var req models.AuthObject

	json.Unmarshal(u.Ctx.Input.RequestBody, &req)

	if (req.Login == "" || req.Password == "") {
		error := struct {
			ErrLoginRequired    bool `json:"errLoginRequired"`
			ErrPasswordRequired bool `json:"errPasswordRequired"`
		}{
			ErrLoginRequired:    req.Login == "",
			ErrPasswordRequired: req.Password == "",
		}

		var jsonError, _= json.Marshal(error);
		u.Data["json"] = customErrors.FormatError(&customErrors.RespError{Code: 500, Message: "Validation failed", Data: string(jsonError)})
		u.ServeJSON()

		return
	}

	token, err := models.GetToken(req.Login, req.Password)

	if err != nil {
		u.Data["json"] = customErrors.FormatError(err)
	} else {
		u.Data["json"] = map[string]string{"accessToken": token["accessToken"], "refreshToken": token["refreshToken"]}
	}

	u.ServeJSON()
}

// @router /refresh-token [post]
func (u *UserController) RefreshToken() {
	var req struct {
		RefreshToken string
	}

	userIdForRefresh := u.Ctx.Input.Session("userIdForRefresh")

	if userIdForRefresh == nil {
		userIdForRefresh = ""
	}

	userId := userIdForRefresh.(string)

	if userId == "" {
		u.Data["json"] = customErrors.FormatError(&customErrors.RespError{Code: 404, Message: "User not found"})
		u.ServeJSON()
		return
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &req)

	token, err := models.RefreshToken(userId, req.RefreshToken)

	if err != nil {
		u.Data["json"] = customErrors.FormatError(err)
	} else {
		u.Data["json"] = token
	}

	u.ServeJSON()
}
