package api

import (
	"context"
	// "database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/sava/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)
// type loginRequest struct {
// 	ExternalID sql.NullInt64 `json:"external_id"`
// 	Name       string        `json:"name"`
// 	Email      string        `json:"email"`
// }

// OAuth config for browser login (only used for "login via browser" testing).
func oauthConfig() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     env.Getenv("OIDC_CLIENT_ID","my-client-id"),
		ClientSecret: env.Getenv("OIDC_CLIENT_SECRET","myclientSecret"), // optional for some providers
		RedirectURL:  env.Getenv("OIDC_REDIRECT_URL","myRedirecturl"),  // e.g., http://localhost:8080/auth/callback
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		Endpoint:     google.Endpoint, // for Google. For other providers use provider.Endpoint()
	}
    return  config
}

func (server *Server) loginUser(ctx *gin.Context){
conf := oauthConfig()
	state := "state" 
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusFound, url)
}


func (server *Server) CallbackHandler(c *gin.Context) {
	conf := oauthConfig()
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("code missing")))
		return
	}

	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("exchange failed: %w", err)))
		return
	}

	rawIDToken, ok := tok.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("id_token missing")))
		return
	}

	// Parse claims from id_token
	var claims struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Sub   string `json:"sub"`
		Phone string `json:"phone"`
		CreatedAt time.Time `json:"created_at"`
	}

	// NOTE: this only decodes the JWT payload (doesn't verify signature here)
	parts := strings.Split(rawIDToken, ".")
	if len(parts) < 2 {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid id_token")))
		return
	}
	decoded, _ := base64.RawURLEncoding.DecodeString(parts[1])
	_ = json.Unmarshal(decoded, &claims)

	resp := map[string]interface{}{
		"access_token": tok.AccessToken,
		"id_token":     rawIDToken,
		"token_type":   tok.TokenType,
		"expiry":       tok.Expiry,
		"user": gin.H{
			"email": claims.Email,
			"name":  claims.Name,
            "sub": claims.Sub,
            "created_at": claims.CreatedAt,
            "phone": claims.Phone,
		},
	}
	c.JSON(http.StatusOK, resp)
}