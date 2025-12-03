package app

import (
	"backend/config"
	"backend/models"
	"context"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/fvbock/endless"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type ValidationError struct {
	Error string `json:"error"`
	Values map[string]string `json:"values"`
}

func Run(cfg config.Config) {
	router := gin.Default()
	
	// Adicionando conexão com o banco de dados ao contexto
	db, err := models.ConnectDb(cfg)
	if err != nil {
		log.Fatal("error DB: %+V", err)
	}

	// servidor de email
	mailjetClient := mailjet.NewMailjetClient(cfg.MJ_APIKEY_PUBLIC, cfg.MJ_APIKEY_PRIVATE)
	
	router.Use(func (c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "db", db))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "mailjet", mailjetClient))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "email", cfg.Email_Sender))
		c.Next()
	})

	// Declarado rota privada que precisa de autorização
	private := router.Group("/private")
	private.Use(authorization)
	private.POST("/echo", echo)

	// Declarado rotas públicas
	public := router.Group("/public")
	public.POST("/create_user", createUser)
	public.POST("/request_reset_password", RequestResetPassword)
	public.POST("/reset_password", ResetPassword)

	// Iniciando o servidor
	endless.ListenAndServe(":8080", router)
}

func echo(c *gin.Context) {
	var body []byte
	buffer := make([]byte, 4)

	for {
		num, err := c.Request.Body.Read(buffer)

		if num > 0 {
			body = append(body, buffer[:num]...)
		}

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, "Failed to read request body")
			return
		}
	}

	c.JSON(http.StatusOK, string(body))

}

func SendBindError(c *gin.Context, err error) {
	var valErr validator.ValidationErrors
	
	// Adicionando o erro ao contexto
	c.Error(err)
	
	// Tentado converter o error em um validation errors
	// Se não conseguir retorna um erro generico
	ok := errors.As(err, &valErr)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extrair os erros de validação, adiciona no mapa 
	values := make(map[string]string)
	for _, value := range valErr {
		values[value.StructField()] = value.Error()
	}

	// retorna o erro
	c.JSON(http.StatusBadRequest, ValidationError {
		Error: "validation error",
		Values: values,
	})

}
