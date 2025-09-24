package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/sava/db/sqlc"
)

// CreateOrderRequest payload
type CreateOrderRequest struct {
    Items []db.OrderItemInput `json:"items"`
}

// CreateOrderHandler handles order creation
func (server *Server) CreateOrderHandler() gin.HandlerFunc {
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
		order, total, err := server.store.CreateOrderWithItems(ctx, customerID, req.Items)
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



// GET /orders/:id
func (server *Server) GetOrderByIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
			return
		}

		order, err := server.store.GetOrderByID(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

// GET /customers/:customer_id/orders
func (server *Server) GetOrdersByCustomerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		custStr := c.Param("customer_id")
		custID, err := strconv.ParseInt(custStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
			return
		}

		orders, err := server.store.GetOrdersByCustomerID(c, custID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}
