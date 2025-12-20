package app

import (
	"backend/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createMovie(c *gin.Context) {
	movie := models.Movie{}
	err := c.BindJSON(&movie)
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
	
	err = models.Create(db, &movie)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "failed to create movie")
		return
	}
	
	c.JSON(http.StatusOK, nil)

}
