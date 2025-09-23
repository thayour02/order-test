package api

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sava/db/middle"
	db "github.com/sava/db/sqlc"
)


type Server struct {
	store *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
// init oidc
	if err := middle.InitOIDC(context.Background()); err != nil {
		log.Println("Warning: OIDC not initialized:", err)
		// if you don't set OIDC envs, login will fail; you can still use other endpoints
	}

	server := &Server{store: store}

	router := gin.Default()

	router.GET("/login", server.loginUser)
	router.GET("/auth/callback", CallbackHandler)

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}