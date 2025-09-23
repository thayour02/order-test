package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sava/utils"
	"github.com/stretchr/testify/require"
)
func createRandomCategory(t *testing.T, parentID *int32) Category {
	arg := CreateCategoryParams{
		Name:     utils.RandomName(),
		ParentID: sql.NullInt64{Int64: 0, Valid: false},
	}

	if parentID != nil {
		arg.ParentID = sql.NullInt64{Int64: int64(*parentID), Valid: true}
	}

	cat, err := testQueries.CreateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, cat.ID)
	return cat
}


// helper to create category
func createTestCategory(t *testing.T, parentID sql.NullInt64) Category {
	if parentID.Valid {
		int32ParentID := int32(parentID.Int64)
		return createRandomCategory(t, &int32ParentID)
	}
	return createRandomCategory(t, nil)
}

func TestCreateCategory(t *testing.T) {
	// root category (no parent)
	root := createTestCategory(t, sql.NullInt64{Valid: false})
	require.NotZero(t, root.ID)

	// child category
	child := createTestCategory(t, sql.NullInt64{Int64: root.ID, Valid: true})
	require.NotZero(t, child.ID)
	require.Equal(t, root.ID, child.ParentID.Int64)
}

func TestAvgPriceForCategory(t *testing.T) {
	// 1️⃣ Create root category
	root := createTestCategory(t, sql.NullInt64{Valid: false})

	// 2️⃣ Create product
	prod := createRandomProduct(t) 

	// 3️⃣ Link product to category
	_, err := testQueries.db.ExecContext(context.Background(),
		`INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2)`,
		prod.ID, root.ID)
	require.NoError(t, err)

	// 4️⃣ Call AvgPriceForCategory
	got, err := testQueries.AvgPriceForCategory(context.Background(), root.ID)
	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.Equal(t, prod.ID, got.ID)
	require.Equal(t, prod.Name, got.Name)
}
