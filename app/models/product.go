package models

import (
	"time"
)

//go:generate goqueryset -in product.go

// Product model
// gen:qs
type Product struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id"`
	ProductName  string     `json:"product_name"`
	Desc         string     `json:"desc"`
	Condition    int32      `json:"condition"`
	ConditionAvg float64    `json:"condition_avg"`
	StartPrice   float64    `json:"start_price"`
	BidMultpl    float64    `json:"bid_multpl"`
	ClosedAT     *time.Time `json:"closed_at"`
	CreatedAT    *time.Time `json:"created_at"`
}

// ProductBidder models
// gen:qs
type ProductBidder struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	ProductID int64      `json:"product_id"`
	Winner    bool       `json:"winner"`
	CreatedAT *time.Time `json:"created_at"`
}

// ProductImage model
// gen:qs
type ProductImage struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	ImageURL  string `json:"image_url"`
}

// ProductLabel model
// gen:qs
type ProductLabel struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
}
