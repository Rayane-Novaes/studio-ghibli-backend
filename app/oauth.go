package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	scope "google.golang.org/api/oauth2/v2"
)



func decodeJwtPart(name, data string) (map[string]any, error)  {
	// XXX: for some reason, this isn't properly decoding...
	buf, err := base64.URLEncoding.DecodeString(data)
	//if err != nil {
	//	panic(fmt.Sprintf("%s: %+v", name, err))
	//}

	if len(buf) == 0 || buf[0] != '{' {
		return nil, fmt.Errorf("%s: invalid jwt", name)
		
	}
	if buf[len(buf)-1] != '}' {
		buf = append(buf, '}')
	}

	m := make(map[string]any)
	err = json.Unmarshal([]byte(buf), &m)
	if err != nil {
		return nil, fmt.Errorf("%s: %+v", name, err)
	}

	return m, nil
}

func (r *RouteData) setupOAuth() {
	conf := &oauth2.Config{
		ClientID:     r.google_client_id,
		ClientSecret: r.google_client_secret,
		Scopes:       []string{
			scope.OpenIDScope,
			scope.UserinfoEmailScope,
			scope.UserinfoProfileScope,
		},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/oauth/callback",
	}

	r.oauth2Config = conf
	r.oauth2Verifier = oauth2.GenerateVerifier()

}

func (r RouteData) oauth(c *gin.Context) {
	
	url := r.oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(r.oauth2Verifier))

	c.Redirect(http.StatusSeeOther, url)
}

func (r RouteData) oauthCallback(c *gin.Context) {
		code, ok := c.GetQuery("code")

		if !ok {
			c.JSON(http.StatusBadRequest, ValidationError{
				Error:  "validation error",
				Values: map[string]string{"code": "Field is missing"},
			})
			return
		}

		tok, err := r.oauth2Config.Exchange(c.Request.Context(), code, oauth2.VerifierOption(r.oauth2Verifier))
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest,ValidationError{
				Error:  "validation error",
				Values: map[string]string{"code": "Invalid code"},
			})
			return
		}

		jwt, _ := tok.Extra("id_token").(string)
		parts := strings.Split(jwt, ".")
		if len(parts) != 3 {
			c.JSON(http.StatusInternalServerError, "Invalid JWT")
			return
		}

		payload, err := decodeJwtPart("payload", parts[1])
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, "Failed to decode JWT")
			return
		}

		// todo: create acess token
		// todo: register user
		// payload["email"]
		// payload["name"]
		

		c.JSON(http.StatusOK, payload)
}