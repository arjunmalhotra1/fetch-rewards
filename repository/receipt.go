package repository

type Receipt struct {
	Retailer     string `json:"retailer" validate:"required,retailer"`
	PurchaseDate string `json:"purchaseDate" validate:"required,purchaseDate"`
	PurchaseTime string `json:"purchaseTime" validate:"required,purchaseTime"`
	Items        []Item `json:"items" validate:"required"`
	Total        string `json:"total" validate:"required,total"`
	UUID         string
	Points       int
}

type Item struct {
	ShortDescription string `json:"shortDescription" validate:"required,shortDescription"`
	Price            string `json:"price" validate:"required,price"`
}
