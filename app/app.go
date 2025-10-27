package app

import (
	"backend/config"
	"backend/models"
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
)

type Handler struct {
	customerHandlerPublic  http.Handler
	customerHandlerPrivate *http.ServeMux
}

func Run(cfg config.Config) {
	serveMuxPublic := http.NewServeMux()
	serveMuxPrivate := http.NewServeMux()

	h := Handler{
		customerHandlerPublic:  serveMuxPublic,
		customerHandlerPrivate: serveMuxPrivate,
	}

	db, err := models.ConnectDb(cfg)
	if err != nil {
		log.Fatal("error DB: %+V", err)
	}

	serveMuxPrivate.HandleFunc("/echo", echo)
	serveMuxPublic.HandleFunc("/create_user", createUser)
	http.HandleFunc("/", h.defaultHandler)

	s := http.Server{
		Addr: ":8080",
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(ctx, "db", db)
		},
	}

	log.Fatal(s.ListenAndServe())
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
	_, pattern := h.customerHandlerPrivate.Handler(r)
	if pattern != "" {
		permission := authorization(r, w)
		if permission != true {
			return
		}

		h.customerHandlerPrivate.ServeHTTP(w, r)
		return
	}

	h.customerHandlerPublic.ServeHTTP(w, r)
}
