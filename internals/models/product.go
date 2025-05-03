package models

type ProductInfo struct {
	Title string  `json:"title"`
	Price float64 `json:"price"`
	Url   string  `json:"url"`
}
