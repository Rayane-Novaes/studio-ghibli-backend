package app

import (
	"backend/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateMovie godoc
// @Summary Create Movie
// @Description add new movie
// @Tags movie
// @Produce json
// @Security BasicAuth
// @Param request body models.Movie true "the request body"
// @Success 200 "Success"
// @Failure 400 {object} ValidationError "Bad Request: Some validation error occurred"
// @Failure 500 {string} string "Internal Server Error: Some unexpected error occurred"
// @Router /private/create_movie [post]
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
