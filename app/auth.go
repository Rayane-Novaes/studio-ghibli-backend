package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"backend/models"
)

func (r RouteData) authorization(c *gin.Context) {

	username, password, ok := c.Request.BasicAuth()

	if !ok {
		c.Error(errors.New("error failed usermanem and password"))
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}

	err := models.ValidUser(r.db, username, []byte(password))

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}

	c.Next()

}

