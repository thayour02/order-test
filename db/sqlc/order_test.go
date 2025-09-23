package db

import (
	"context"
	"testing"
	"time"


	"github.com/stretchr/testify/require"
)


func TestCreateOrderAndOrderItem(t *testing.T) {
	customer := createRandomCustomer(t)
	product := createRandomProduct(t)

	// --- Create Order ---
	orderArg := CreateOrderParams{
		CustomerID: int64(customer.ID),
		Total:      "30.00",
	}

	order, err := testQueries.CreateOrder(context.Background(), orderArg)
	require.NoError(t, err)
	require.NotZero(t, order.ID)
	require.Equal(t, orderArg.Total, order.Total)
	require.Equal(t, customer.ID, order.CustomerID)
	require.WithinDuration(t, time.Now(), order.CreatedAt, time.Second*2)

	// --- Create Order Item ---
	itemArg := CreateOrderItemParams{
		OrderID:   int64(order.ID),
		ProductID: int64(product.ID),
		Quantity:  2,
		Subtotal:  "30.00",
	}

	item, err := testQueries.CreateOrderItem(context.Background(), itemArg)
	require.NoError(t, err)
	require.Equal(t, itemArg.OrderID, item.OrderID)
	require.Equal(t, itemArg.ProductID, item.ProductID)
	require.Equal(t, itemArg.Quantity, item.Quantity)
	require.Equal(t, itemArg.Subtotal, item.Subtotal)
}
