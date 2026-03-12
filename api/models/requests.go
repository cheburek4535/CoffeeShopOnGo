// api/models/requests.go
package models

type CreateProductRequest struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Category string  `json:"category" binding:"required"`
	VAT      int     `json:"vat" binding:"min=0,max=100"`
	IsActive bool    `json:"is_active"`
}

type UpdatePriceRequest struct {
	Price float64 `json:"price" binding:"required,gt=0"`
}

type CreateStaffRequest struct {
	FullName string  `json:"full_name" binding:"required"`
	Salary   float64 `json:"salary" binding:"required,gt=0"`
	Age      int     `json:"age" binding:"required,min=18,max=100"`
	Position string  `json:"position" binding:"required"`
}

type CreateSaleRequest struct {
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	ProductID int     `json:"product_id" binding:"required"`
	StaffID   int     `json:"staff_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
}
