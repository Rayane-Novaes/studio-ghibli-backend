package app

import (
	"backend/models"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=12,lte=<=72,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz,contaisany=0123456789,containsany=|@#%$^[]{}?!*~();."`
	Email    string `json:"email" binding:"required,email"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validator.New().Struct(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	db, ok := r.Context().Value("db").(models.DB)
	if ok != true {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = models.CreateUser(db, user.Username, []byte(user.Password))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func authorization(r *http.Request, w http.ResponseWriter) bool {
	username, password, ok := r.BasicAuth()

	if ok != true {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	db, ok := r.Context().Value("db").(models.DB)
	if ok != true {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	err := models.ValidUser(db, username, []byte(password))

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}
