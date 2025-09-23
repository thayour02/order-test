package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sava/utils"
	"github.com/stretchr/testify/require"
)

// helper to create a random OIDC customer
func createRandomCustomer(t *testing.T) Customer {
	arg := CreateCustomerParams{
		OidcSub: utils.RandomString(32), 
		Name:    utils.RandomName(),
		Email:   utils.RandomEmail(),
		Phone:   sql.NullInt64{Int64: 1234567890, Valid: true},
	}

	customer, err := testQueries.CreateCustomer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, customer)

	require.Equal(t, arg.OidcSub, customer.OidcSub)
	require.Equal(t, arg.Name, customer.Name)
	require.Equal(t, arg.Email, customer.Email)
	require.Equal(t, arg.Phone.Int64, customer.Phone.Int64)

	require.NotZero(t, customer.ID)
	require.NotZero(t, customer.CreatedAt)

	return Customer{
		ID:        customer.ID,
		OidcSub:   customer.OidcSub,
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		CreatedAt: customer.CreatedAt,
	}
}

func TestCreateCustomer(t *testing.T) {
	createRandomCustomer(t)
}

func TestGetCustomerByOIDCSub(t *testing.T) {
	customer1 := createRandomCustomer(t)

	customer2, err := testQueries.GetCustomerByOIDCSub(context.Background(), customer1.OidcSub)
	require.NoError(t, err)
	require.NotEmpty(t, customer2)

	require.Equal(t, customer1.ID, customer2.ID)
	require.Equal(t, customer1.OidcSub, customer2.OidcSub)
	require.Equal(t, customer1.Name, customer2.Name)
	require.Equal(t, customer1.Email, customer2.Email)
	require.Equal(t, customer1.Phone, customer2.Phone)
}
