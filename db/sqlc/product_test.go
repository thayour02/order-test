package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/sava/utils"
	"github.com/stretchr/testify/require"
)


func createRandomProduct(t *testing.T) Product {
arg := CreateProductParams{
	Name:        utils.RandomName(),
	Description: sql.NullString{String: utils.RandomDescription(), Valid: true},
	Price:       utils.RandomPrice(),
}

	product, err := testQueries.CreateProduct(context.Background(), arg)
if err != nil {
    t.Errorf("CreateProduct failed: %v", err)
}
require.NoError(t, err)

	require.NotEmpty(t, product)


	require.Equal(t, arg.Name, product.Name)
	require.Equal(t, arg.Description, product.Description)
	require.Equal(t, arg.Price, product.Price)
	return product
}

func TestCreateProduct(t *testing.T){
 	createRandomProduct(t)
}

func TestCreateAndGetProduct(t *testing.T) {
	prod := createRandomProduct(t)

	got, err := testQueries.GetProductByID(context.Background(), int64(prod.ID))
	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.Equal(t, prod.ID, got.ID)
	require.Equal(t, prod.Name, got.Name)
	require.Equal(t, prod.Description, got.Description)
	require.Equal(t, prod.Price, got.Price)
	require.WithinDuration(t, prod.CreatedAt, got.CreatedAt, time.Second)
}

func TestAddProductCategory(t *testing.T) {
	prod := createRandomProduct(t)
	cat := createRandomCategory(t, nil)

	arg := AddProductCategoryParams{
		ProductID:  int64(prod.ID),
		CategoryID: int64(cat.ID),
	}

	err := testQueries.AddProductCategory(context.Background(), arg)
	require.NoError(t, err)

	// check recursive query
	products, err := testQueries.ProductsInCategoryRecursive(context.Background(), int64(cat.ID))
	require.NoError(t, err)
	require.NotEmpty(t, products)

	found := false
	for _, p := range products {
		if p.ID == prod.ID {
			found = true
			break
		}
	}
	require.True(t, found)
}
