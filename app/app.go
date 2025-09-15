package app

import (
	"backend/config"
	"backend/models"
	"errors"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	customerHandler http.Handler
}

func Run(cfg config.Config) {
	serveMux := http.NewServeMux()
	h := Handler{
		customerHandler: serveMux,
	}

	_, err := models.ConnectDb(cfg)
	if err != nil {
		log.Fatal("error DB: %+V", err)
	}

	serveMux.HandleFunc("/echo", echo)
	http.HandleFunc("/", h.defaultHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	var body []byte
	buffer := make([]byte, 4)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	for {
		num, err := r.Body.Read(buffer)

		if num > 0 {
			body = append(body, buffer[:num]...)
		}

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func (h *Handler) defaultHandler(w http.ResponseWriter, r *http.Request) {
	h.customerHandler.ServeHTTP(w, r)
}
