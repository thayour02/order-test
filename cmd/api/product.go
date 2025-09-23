package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/sava/db/sqlc"
)

// CreateProductRequest is the API payload
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       string  `json:"price" binding:"required"`
	CategoryIDs []int64 `json:"category_ids"`
}

// CreateProductHandler handles product creation
func CreateProductHandler(store *db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateProductRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product, err := store.CreateProductWithCategories(ctx, db.ProductInput{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			CategoryIDs: req.CategoryIDs,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description.String,
			"price":       product.Price,
		})
	}
}
