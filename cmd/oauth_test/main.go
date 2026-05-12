package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	scope "google.golang.org/api/oauth2/v2"
)

func decodeJwtPart(name, data string) map[string]any {
	// XXX: for some reason, this isn't properly decoding...
	buf, err := base64.URLEncoding.DecodeString(data)
	//if err != nil {
	//	panic(fmt.Sprintf("%s: %+v", name, err))
	//}

	if len(buf) == 0 || buf[0] != '{' {
		panic(fmt.Sprintf("%s: invalid jwt", name))
	}
	if buf[len(buf)-1] != '}' {
		buf = append(buf, '}')
	}

	m := make(map[string]any)
	err = json.Unmarshal([]byte(buf), &m)
	if err != nil {
		panic(fmt.Sprintf("%s: %+v", name, err))
	}

	fmt.Printf("%s: %+v\n", name, m)
	return m
}

func main() {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{
			scope.OpenIDScope,
			scope.UserinfoEmailScope,
			scope.UserinfoProfileScope,
		},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/google/callback",
	}

	// use PKCE to protect against CSRF attacks
	// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
	verifier := oauth2.GenerateVerifier()

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))

			fmt.Printf("%+v\n", err)
			return
		}

		jwt, _ := tok.Extra("id_token").(string)
		parts := strings.Split(jwt, ".")
		if len(parts) != 3 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))

			fmt.Printf("bad jwt: '%s'\n", jwt)

			return
		}

		_ = decodeJwtPart("header", parts[0])
		_ = decodeJwtPart("payload", parts[1])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusSeeOther)
	})

	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		srv.Shutdown(context.Background())
		srv.Close()
	}()

	fmt.Printf("waiting...\n")

	err := srv.ListenAndServe()
	fmt.Printf("exit err: %+v\n", err)
}