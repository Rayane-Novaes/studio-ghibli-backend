package app

import (
	"backend/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

var router http.Handler

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load envs")
	}

	router = setup(cfg)

	m.Run()
}

func SendRequest(t *testing.T, method, target string, header map[string]string, data any) *http.Response {

	var body io.ReadCloser

	if data != nil {
		payload, err := json.Marshal(data)
		require.NoError(t, err, "failed json marshal")

		reader := bytes.NewBuffer(payload)
		body = io.NopCloser(reader)

	}

	req := httptest.NewRequest(method, target, body)
	w := httptest.NewRecorder()

	for key, value := range header {
		req.Header.Add(key, value)
	}

	router.ServeHTTP(w, req)

	return w.Result()
}

// TODO: verificar campos retornados e finalizar os testes
func TestCreateUser(t *testing.T) {
	user := User{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 12),
	}

	response := SendRequest(t, "POST", "/public/create_user", nil, user)

	// Comparações
	require.Equal(t, 200, response.StatusCode)

}

