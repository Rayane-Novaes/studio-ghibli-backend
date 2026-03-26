package app

import (
	"backend/models"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestListMoviesEmpty testa se a listagem de filmes retorna status 200
func TestListMoviesEmpty(t *testing.T) {
	response := SendRequest(t, "GET", "/public/list_movies", nil, nil)

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var result PaginationReturn
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	require.NotNil(t, result.Data)
}

// TestListMoviesWithCursor testa a listagem de filmes com cursor
func TestListMoviesWithCursor(t *testing.T) {
	response := SendRequest(t, "GET", "/public/list_movies?cursor=", nil, nil)

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var result PaginationReturn
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	require.NotNil(t, result.Data)
}

// TestCreateMovieSuccess testa a criação de um filme com dados válidos
func TestCreateMovieSuccess(t *testing.T) {
	// Criar uma imagem base64 válida (PNG 1x1 pixel)
	validBase64Image := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

	movie := models.Movie{
		Name:         "Spirited Away",
		Director:     "Hayao Miyazaki",
		Producer:     "Toshio Suzuki",
		ReleaseDate:  "2001-07-20",
		Duration:     models.Duration(125 * 60 * 1e9),
		BannerImagem: validBase64Image,
	}

	// Criar header com autenticação básica
	headers := map[string]string{
		"Authorization": "Basic dXNlcjpwYXNz", // user:pass em base64
	}

	response := SendRequest(t, "POST", "/private/create_movie", headers, movie)

	require.Equal(t, http.StatusNoContent, response.StatusCode)
}

// TestCreateMovieInvalidImage testa a criação de um filme com imagem inválida
func TestCreateMovieInvalidImage(t *testing.T) {
	// String base64 inválida que não é uma imagem
	invalidBase64Image := base64.StdEncoding.EncodeToString([]byte("not an image"))

	movie := models.Movie{
		Name:         "Howl's Moving Castle",
		Director:     "Hayao Miyazaki",
		Producer:     "Toshio Suzuki",
		ReleaseDate:  "2004-11-20",
		Duration:     models.Duration(119 * 60 * 1e9),
		BannerImagem: invalidBase64Image,
	}

	headers := map[string]string{
		"Authorization": "Basic dXNlcjpwYXNz",
	}

	response := SendRequest(t, "POST", "/private/create_movie", headers, movie)

	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var result ValidationError
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	require.Equal(t, "validation error", result.Error)
}

// TestCreateMovieMissingFields testa a criação de um filme sem campos obrigatórios
func TestCreateMovieMissingFields(t *testing.T) {
	// Movie sem o campo Name (obrigatório)
	incompleteMovie := map[string]interface{}{
		"director":     "Hayao Miyazaki",
		"producer":     "Toshio Suzuki",
		"release_date": "2001-07-20",
		"duration":     "2h5m",
		"banner_image": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
	}

	headers := map[string]string{
		"Authorization": "Basic dXNlcjpwYXNz",
	}

	response := SendRequest(t, "POST", "/private/create_movie", headers, incompleteMovie)

	require.Equal(t, http.StatusBadRequest, response.StatusCode)
}
