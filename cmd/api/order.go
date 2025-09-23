package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/sava/db/sqlc"
)

// CreateOrderRequest payload
type CreateOrderRequest struct {
    Items []db.OrderItemInput `json:"items"`
}

// CreateOrderHandler handles order creation
func CreateOrderHandler(store *db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateOrderRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(req.Items) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "order must contain at least one item"})
			return
		}

		// Get customer ID from context
		val, ok := ctx.Get("customer_id")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "customer_id not found in context"})
			return
		}
		customerID := val.(int64)

		// Create order with store
		order, total, err := store.CreateOrderWithItems(ctx, customerID, req.Items)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"order_id": order.ID,
			"total":    fmt.Sprintf("%.2f", total),
		})
	}
}
