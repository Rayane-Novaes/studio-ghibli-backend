package app

import (
	"backend/models"
	"bytes"
	"encoding/base64"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaginationReturn struct {
	Data   any
	Cursor string
}

type ListerParam struct {
	Cursor string `url:"cursor"`
}

// CreateMovie godoc
// @Summary Create Movie
// @Description add new movie
// @Tags movie
// @Produce json
// @Security BasicAuth
// @Param request body models.Movie true "the request body"
// @Success 204 "Success"
// @Failure 400 {object} ValidationError "Bad Request: Some validation error occurred"
// @Failure 500 {string} string "Internal Server Error: Some unexpected error occurred"
// @Router /private/create_movie [post]
func (r RouteData) createMovie(c *gin.Context) {
	movie := models.Movie{}
	err := c.BindJSON(&movie)
	if err != nil {
		SendBindError(c, err)
		return
	}

	if !r.skip_image {
		err = validationImage([]byte(movie.BannerImagem))
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, ValidationError{
				Error:  "validation error",
				Values: map[string]string{"banner_image": "Invalid image: must be a base64 encoded png, jpg or gif"},
			})
			return
		}
	}

	err = models.Create(r.db, &movie)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "failed to create movie")
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

func validationImage(data []byte) error {
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(data))
	_, _, err := image.Decode(reader)

	return err

}

// TODO: Adicionar a documentação swagger
func (r RouteData) listMovie(c *gin.Context) {
	lister := []*models.Movie{}
	lister_param := ListerParam{}
	err := c.BindQuery(&lister_param)
	if err != nil {
		SendBindError(c, err)
		return
	}

	cursor, err := models.List(r.db, &lister, lister_param.Cursor)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "failed list movies")
		return
	}

	c.JSON(http.StatusOK, &PaginationReturn{Data: lister, Cursor: cursor})

}
