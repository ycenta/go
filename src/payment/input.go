package payment

type InputPayment struct {
	ProductID int     `json:"productid" binding:"required"`
	PricePaid float64 `json:"pricepaid" binding:"required"`
}
