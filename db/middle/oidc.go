package middle

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	db "github.com/sava/db/sqlc"
	"github.com/sava/env"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// OIDCVerifier holds provider/verifier for reuse
type OIDCVerifier struct {
	Provider *oidc.Provider
	Verifier *oidc.IDTokenVerifier
}

var verifier *OIDCVerifier

// InitOIDC initializes the OIDC provider and verifier
func InitOIDC(ctx context.Context) error {
	issuer := env.Getenv("OIDC_ISSUER", "https://accounts.google.com")
	clientID := env.Getenv("OIDC_CLIENT_ID", "")
	if clientID == "" {
		return fmt.Errorf("OIDC_CLIENT_ID must be set")
	}
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return err
	}
	verifier = &OIDCVerifier{
		Provider: provider,
		Verifier: provider.Verifier(&oidc.Config{ClientID: clientID}),
	}
	return nil
}

// AuthMiddleware verifies bearer token and ensures customer exists
func  AuthMiddleware(store *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "bearer token is required"})
			c.Abort()
			return
		}
		raw := strings.TrimSpace(auth[len("Bearer "):])

		idtoken, err := verifier.Verifier.Verify(c.Request.Context(), raw)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			c.Abort()
			return
		}

		var claims struct {
			OidcSub string `json:"sub"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Phone   string `json:"phone"`
		}
		if err := idtoken.Claims(&claims); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to parse token claims"})
			c.Abort()
			return
		}

		if claims.Email == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token missing email"})
			c.Abort()
			return
		}

		// use store (not q)
		customer, err := store.GetCustomerByOIDCSub(c, claims.OidcSub)
		if err != nil {
			createdCustomer, err := store.CreateCustomer(c, db.CreateCustomerParams{
				OidcSub: claims.OidcSub,
				Name:    claims.Name,
				Email:   claims.Email,
				Phone:   sql.NullInt64{Valid: false}, // or fill if you have phone
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create customer"})
				c.Abort()
				return
			}
			// Map db.Customer to db.GetCustomerByOIDCSubRow
			customer = db.GetCustomerByOIDCSubRow{
				ID:      createdCustomer.ID,
				OidcSub: createdCustomer.OidcSub,
				Name:    createdCustomer.Name,
				Email:   createdCustomer.Email,
				Phone:   createdCustomer.Phone,
				// Add other fields if needed
			}
		}
		c.Set("customer", customer)
		c.Next()
	}
}
