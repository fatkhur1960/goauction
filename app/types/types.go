package types

import (
	"time"
)

type (
	// Product definisi type untuk menampilkan product
	Product struct {
		ID            int64       `json:"id"`
		ProductName   string      `json:"product_name"`
		ProductImages interface{} `json:"product_images"`
		Desc          string      `json:"desc"`
		Condition     int32       `json:"condition" `
		ConditionAvg  float64     `json:"condition_avg"`
		StartPrice    float64     `json:"start_price"`
		BidMultpl     float64     `json:"bid_multpl"`
		ClosedAT      *time.Time  `json:"closed_at"`
		CreatedAT     *time.Time  `json:"created_at"`
		Labels        interface{} `json:"labels"`
		Sold          bool        `json:"sold"`
		Closed        bool        `json:"closed"`
		BidStatus     interface{} `json:"bid_status"`
	}

	// ProductDetail json
	ProductDetail struct {
		ID            int64       `json:"id"`
		ProductName   string      `json:"product_name"`
		ProductImages interface{} `json:"product_images"`
		Desc          string      `json:"desc"`
		Condition     int32       `json:"condition" `
		ConditionAvg  float64     `json:"condition_avg"`
		StartPrice    float64     `json:"start_price"`
		BidMultpl     float64     `json:"bid_multpl"`
		ClosedAT      *time.Time  `json:"closed_at"`
		CreatedAT     *time.Time  `json:"created_at"`
		Labels        interface{} `json:"labels"`
		Sold          bool        `json:"sold"`
		Closed        bool        `json:"closed"`
		Bidders       interface{} `json:"bidders"`
	}

	// Chat api type
	Chat struct {
		ID           int64       `json:"id"`
		InitiatorID  int64       `json:"initiator_id"`
		SubscriberID int64       `json:"subscriber_id"`
		LastUpdated  *time.Time  `json:"last_updated"`
		TS           *time.Time  `json:"ts"`
		Display      interface{} `json:"display"`
	}
)
