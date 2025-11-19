package app

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/sirupsen/logrus"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=12,lte=<=72"`
	Email    string `json:"email" binding:"required,email"`
}

type RequestResetPasswordBody struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordBody struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,gte=12,lte=<=72"`
	Password    string `json:"password" binding:"required,password"`
}

func createUser(c *gin.Context) {
	user := User{}
	err := c.BindJSON(&user)
	if err != nil {
		SendBindError(c, err)
		return
	}
	
	validatorPassword := ValidatorPassword(user.Password)
	if !validatorPassword {
		c.JSON(http.StatusBadRequest, ValidationError {
			Error: "Validation Errors",
			Values: map[string]string{
				"password": "password must have at least: [A-Z], [a-z], [0-9], [|@#%$^[]{}?!*~();.]",
			},
		})
		return
	}

	db, ok := c.Request.Context().Value("db").(models.DB)
	if ok != true {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	err = models.CreateUser(db, user.Username, []byte(user.Password), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func authorization(c *gin.Context){

	username, password, ok := c.Request.BasicAuth()

	if ok != true {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	db, ok := c.Request.Context().Value("db").(models.DB)
	if ok != true {
		c.JSON(http.StatusInternalServerError, nil)
		c.Abort()
		return 
	}

	err := models.ValidUser(db, username, []byte(password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	c.Next()

}

func RequestResetPassword(c *gin.Context) {
	email := RequestResetPasswordBody{}
	err := json.NewDecoder(c.Request.Body).Decode(&email)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = validator.New().Struct(&email)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	db, ok := c.Request.Context().Value("db").(models.DB)
	if ok != true {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	pwd, err := models.UpdateUserPasswordReset(db, email.Email)
	if err == nil {
		// TODO: enviar email
	}

	c.JSON(http.StatusOK, pwd)
}

func ResetPassword(c *gin.Context) {
	reset := ResetPasswordBody{}
	err := c.BindJSON(&reset)
	if err != nil {
		SendBindError(c, err)
		return
	}

	validatorPassword := ValidatorPassword(reset.NewPassword)
	if !validatorPassword {
		c.JSON(http.StatusBadRequest, ValidationError {
			Error: "Validation Errors",
			Values: map[string]string{
				"new_password": "password must have at least: [A-Z], [a-z], [0-9], [|@#%$^[]{}?!*~();.]",
			},
		})
		return
	}
	
	db, ok := c.Request.Context().Value("db").(models.DB)
	if ok != true {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	err = models.ResetPassword(db, reset.Email, reset.NewPassword, reset.Password)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func ValidatorPassword (password string) (bool) {
	return strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") &&
		strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") &&
		strings.ContainsAny(password, "|@#%$^[]{}?!*~();.") &&
		strings.ContainsAny(password, "0123456789")
}
/*
	Requisição do reset: input (email). Gerar um código unico por usuário
	Reset
*/
