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

	"github.com/gin-gonic/gin"
)

func Run(cfg config.Config) {
	router := gin.Default()
	
	// Adicionando conexão com o banco de dados ao contexto
	db, err := models.ConnectDb(cfg)
	if err != nil {
		log.Fatal("error DB: %+V", err)
	}

	router.Use(func (c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "db", db))
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
