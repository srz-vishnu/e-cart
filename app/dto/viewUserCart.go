package dto

type ViewCart struct {
	ProductID   int64   `json:"product_id"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
	BrandName   string  `json:"brandname"`
	TotalAmount float64 `json:"totalamount"`
}
