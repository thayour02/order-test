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
for _, ri := range router.Routes() {
    log.Printf("Route registered: %s %s", ri.Method, ri.Path)
}

	router.GET("/auth/login", server.loginUser)
	router.GET("/auth/callback", server.CallbackHandler)
	router.GET("/products", server.GetProductsHandler())

	secure := router.Group("/").Use(middle.AuthMiddleware(server.store))
	
	secure.POST("/create-product", server.CreateProductHandler())
	secure.GET("/categories/:id/avg_price", server.AvgPriceHandler())
	secure.POST("/orders", server.CreateOrderHandler())
	secure.GET("/orders/:id", server.GetOrderByIDHandler())

	server.router = router

	return server
	
	
}

func (server *Server) Start(address string) error {
	for _, ri := range server.router.Routes() {
		log.Println("Route:", ri.Method, ri.Path)
	}

	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}