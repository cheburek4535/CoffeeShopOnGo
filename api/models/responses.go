package models

import (
    "time"
)

type ProductResponse struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    Category string  `json:"category"`
    VAT      int     `json:"vat"`
}

type StaffResponse struct {
    ID       int     `json:"id"`
    FullName string  `json:"full_name"`
    Salary   float64 `json:"salary"`
    Position string  `json:"position"`
    Age      int     `json:"age"`
}

type SalesResponse struct {
	ID int `json:"id"`
	Amount float64 `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
    ProductID int       `json:"product_id"`
    StaffID   int       `json:"staff_id"`
    Quantity  int       `json:"quantity"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}