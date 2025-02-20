package models

type CartItem struct {
	UserID int `json:"user_id"`
    ProductID int `json:"product_id"`
    Quantity int `json:"quantity"`
    Price float64 `json:"price"`
}

type Cart struct {
    Items []CartItem `json:"items"`
    Total float64 `json:"total"`
}
