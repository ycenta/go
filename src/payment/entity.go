package payment

import (
	"go/src/product"
	"time"
)

type Payment struct {
	ID        int              `json:"id"`
	ProductID int              `json:"product_id"`
	Product   *product.Product `json:"product"`
	PricePaid float64          `json:"price_paid"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
