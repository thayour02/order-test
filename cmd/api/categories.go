package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/sava/db/sqlc"
)

type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *int32 `json:"parent_id"`
}

type AvgPriceRequest struct {
	CategoryID int64 `json:"category_id" binding:"required"`
}


// CreateCategoryHandler creates a category
func CreateCategoryHandler(store *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		parentID := sql.NullInt64{Valid: false}
		if req.ParentID != nil {
			parentID = sql.NullInt64{Int64: int64(*req.ParentID), Valid: true}
		}

		category, err := store.CreateCategory(c, req.Name, parentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, category)
	}
}



// AvgPriceHandler returns the average product price for a given category

func (server *Server) AvgPriceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryIDParam := c.Param("id")
		fmt.Printf("Received category ID param: %s\n", categoryIDParam)
		categoryID, err := strconv.ParseInt(categoryIDParam, 10, 64)
		fmt.Printf("Parsed category ID: %d\n", categoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
			return
		}

		avgPrice, err := server.store.AvgPriceForCategory(c, categoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"avg_price": avgPrice})
	}
}
