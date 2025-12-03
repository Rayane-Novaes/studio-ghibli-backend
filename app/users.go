package app

import (
	"backend/models"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=12,lte=72"`
	Email    string `json:"email" binding:"required,email"`
}

type RequestResetPasswordBody struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordBody struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,gte=12,lte=72"`
	Password    string `json:"password" binding:"required"`
}

func createUser(c *gin.Context) {
	// Lendo a requisição e validando
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
		c.Error(errors.New("failed db connection"))
		c.JSON(http.StatusInternalServerError, "failed to create user")
		return
	}

	err = models.CreateUser(db, user.Username, []byte(user.Password), user.Email)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "failed to create user")
		return
	}

	c.JSON(http.StatusOK, nil)
}

func authorization(c *gin.Context){

	username, password, ok := c.Request.BasicAuth()

	if !ok {
		c.Error(errors.New("error failed usermanem and password"))
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}

	db, ok := c.Request.Context().Value("db").(models.DB)
	if !ok {
		c.Error(errors.New("failed db connection"))
		c.JSON(http.StatusInternalServerError, "Unauthorized")
		c.Abort()
		return 
	}

	err := models.ValidUser(db, username, []byte(password))

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}

	c.Next()

}

func RequestResetPassword(c *gin.Context) {
	email := RequestResetPasswordBody{}	
	err := c.BindJSON(&email)
	if err != nil {
		SendBindError(c, err)
		return
	}

	db, ok := c.Request.Context().Value("db").(models.DB)
	if ok != true {
		c.Error(errors.New("Failed connect db"))
		c.JSON(http.StatusInternalServerError, "failed to request password reset")
		return
	}

	pwd, username, err := models.UpdateUserPasswordReset(db, email.Email)
	if err == nil {
		err = SendEmail(c, email.Email, username, pwd)
		if err != nil {
			c.Error(err)
		}
	}

	c.JSON(http.StatusOK, pwd)
}

func SendEmail (c *gin.Context, email string, username string, passwordTemp string) (error) {
	m, ok := c.Request.Context().Value("mailjet").(*mailjet.Client)
	if !ok {
		return errors.New("failed mailjet client")
	}

	email_sender, ok := c.Request.Context().Value("email").(string)
	if !ok {
		return errors.New("failed email sender")
	}

	link := "http:localhost:3000/reset_password?email=" + url.QueryEscape(email) + "&temp_password=" + url.QueryEscape(passwordTemp) 

	messagesInfo := []mailjet.InfoMessagesV31 {
      mailjet.InfoMessagesV31{
        From: &mailjet.RecipientV31{
          Email: email_sender,
          Name: "Teste",
        },
        To: &mailjet.RecipientsV31{
          mailjet.RecipientV31 {
            Email: email,
            Name: username,
          },
        },
        Subject: "Reset password",
        TextPart: "We just received a new reset password request for this account. If you don't regonize this request, ignore this email.\n\nOtherwise, click on this link to reset your password:" + link,
      },
    }

	messages := mailjet.MessagesV31{Info: messagesInfo }
	res, err := m.SendMailV31(&messages)
	if err != nil {
		return err
	}
	fmt.Printf("Data: %+v\n", res)
	return nil
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
		c.Error(errors.New("Failed connect db"))
		c.JSON(http.StatusInternalServerError, "failed to reset password")
		return
	}

	err = models.ResetPassword(db, reset.Email, reset.NewPassword, reset.Password)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "failed to reset password")
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
