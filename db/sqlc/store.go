package db

// store provides all functions to execute db queries and transactions

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

)

type Store struct {
	queries *Queries
	db      *sql.DB
}

type OrderItemInput struct {
    ProductID int64
    Quantity  int32
}

type ProductInput struct {
	Name        string
	Description string
	Price       string
	CategoryIDs []int64
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		queries: New(db),
	}

}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := New(tx)
	err = fn(query)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}


// CreateOrderWithItems creates an order and its items inside a transaction
func (store *Store) CreateOrderWithItems(ctx context.Context, customerID int64, items []OrderItemInput) (Order, float64, error) {
	var order Order
	var total float64

	err := store.execTx(ctx, func(q *Queries) error {
		// Create order
		ord, err := q.CreateOrder(ctx, CreateOrderParams{
			CustomerID: customerID,
			Total:      "0.00",
			// CreatedAt:  time.Now(),
		})
		if err != nil {
			return err
		}
		order := ord

		//Create order items and compute total
		for _, it := range items {
			product, err := q.GetProductByID(ctx, it.ProductID)
			if err != nil {
				return err
			}
			price, err := strconv.ParseFloat(product.Price, 64)
			if err != nil {
				return err
			}
			subtotal := price * float64(it.Quantity)
			total += subtotal

			_, err = q.CreateOrderItem(ctx, CreateOrderItemParams{
				OrderID:   order.ID,
				ProductID: it.ProductID,
				Quantity:  it.Quantity,
				Subtotal:  fmt.Sprintf("%.2f", subtotal),
			})
			if err != nil {
				return err
			}
		}

		// 3️⃣ Update order total
		return q.UpdateOrderTotal(ctx, UpdateOrderTotalParams{
			ID:    order.ID,
			Total: fmt.Sprintf("%.2f", total),
		})
	})

	return order, total, err
}


// GetCustomerByOIDCSub retrieves a customer by OIDC sub
func (store *Store) GetCustomerByOIDCSub(ctx context.Context, oidcSub string) (GetCustomerByOIDCSubRow, error) {
	return store.queries.GetCustomerByOIDCSub(ctx, oidcSub)
}

// CreateCustomer wraps the sqlc-generated CreateCustomer query in a transaction
func (store *Store) CreateCustomer(ctx context.Context, arg CreateCustomerParams) (Customer, error) {
	var customer Customer
	err := store.execTx(ctx, func(q *Queries) error {
		row, err := q.CreateCustomer(ctx, arg) // returns CreateCustomerRow
		if err != nil {
			return err
		}
		// map manually
		customer = Customer{
			ID:        row.ID,
			Name:      row.Name,
			Email:     row.Email,
			OidcSub:   row.OidcSub,
			CreatedAt: row.CreatedAt,
			Phone: row.Phone,
		}
		return nil
	})
	return customer, err
}


func (store *Store) CreateProductWithCategories(ctx context.Context, input ProductInput) (Product, error) {
	var product Product

	err := store.execTx(ctx, func(q *Queries) error {
		desc := sql.NullString{String: input.Description, Valid: input.Description != ""}
		p, err := q.CreateProduct(ctx, CreateProductParams{
			Name:        input.Name,
			Description: desc,
			Price:       input.Price,
		})
		if err != nil {
			return err
		}
		product = p

		for _, catID := range input.CategoryIDs {
			err := q.AddProductCategory(ctx, AddProductCategoryParams{
				ProductID:  product.ID,
				CategoryID: catID,
			})
			if err != nil {
				return fmt.Errorf("failed to add product to category %d: %w", catID, err)
			}
		}

		return nil
	})

	return product, err
}

// GetAllProducts wraps Queries.GetAllProducts
func (store *Store) GetAllProducts(ctx context.Context) ([]Product, error) {
	return store.queries.GetAllProducts(ctx)
}


