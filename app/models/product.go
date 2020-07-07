package models

import (
	"time"

	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/types"
)

//go:generate goqueryset -in product.go

// Product model
// gen:qs
type Product struct {
	ID            int64          `json:"id"`
	StoreID       int64          `json:"store_id"`
	ProductImages []ProductImage `json:"product_images" gorm:"foreignkey:ProductID"`
	ProductName   string         `json:"product_name"`
	Desc          string         `json:"desc"`
	Condition     int32          `json:"condition"`
	ConditionAvg  float64        `json:"condition_avg"`
	StartPrice    float64        `json:"start_price"`
	BidMultpl     float64        `json:"bid_multpl"`
	ClosedAT      *time.Time     `json:"closed_at"`
	CreatedAT     *time.Time     `json:"created_at"`
	Sold          bool           `json:"sold"`
	Closed        bool           `json:"closed"`
	Labels        []ProductLabel `json:"labels" gorm:"foreignkey:ProductID"`
}

// ProductBidder models
// gen:qs
type ProductBidder struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	ProductID int64       `json:"product_id"`
	BidPrice  float64     `json:"bid_price"`
	Winner    bool        `json:"winner"`
	CreatedAT *time.Time  `json:"created_at"`
	User      *UserSimple `json:"user,omitempty"`
}

// ProductImage model
// gen:qs
type ProductImage struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"-"`
	ImageURL  string `json:"image_url"`
}

// ProductLabel model
// gen:qs
type ProductLabel struct {
	ID        int64  `json:"-"`
	ProductID int64  `json:"-"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

// BidStatus status bid for product
type BidStatus struct {
	BidCount       int     `json:"bid_count"`
	LatestBidPrice float64 `json:"latest_bid_price"`
	LatestUserID   int64   `json:"latest_user_id"`
	MyLatestBid    float64 `json:"my_latest_bid,omitempty"`
}

// TableName override table name
func (BidStatus) TableName() string {
	return "product_bidders"
}

// GetBidderStatus digunakan untuk mendapatkan status bid product
func (p *Product) GetBidderStatus(userID *int64) BidStatus {
	bidStatus := BidStatus{}
	bidPrice := app.DB.Model(&ProductBidder{}).Select("bid_price").Where("user_id = ?", userID).Order("id DESC").Limit(1)
	latestUserID := app.DB.Model(&ProductBidder{}).Select("user_id").Where("product_id = ?", p.ID).Order("id DESC").Limit(1)
	app.DB.Select(
		"MAX(bid_price) AS latest_bid_price, COUNT(id) AS bid_count, ? AS latest_user_id, ? AS my_latest_bid",
		latestUserID.SubQuery(),
		bidPrice.SubQuery(),
	).Where("product_id = ?", p.ID).Find(&bidStatus)

	return bidStatus
}

// GetLatestBidPrice from product
func (p *Product) GetLatestBidPrice() float64 {
	bidStatus := BidStatus{}
	app.DB.Select("MAX(bid_price) AS latest_bid_price").Where("product_id = ?", p.ID).First(&bidStatus)

	return bidStatus.LatestBidPrice
}

// ToAPI --
func (p *Product) ToAPI(userID *int64) types.Product {

	images := []ProductImage{}
	labels := []ProductLabel{}
	app.DB.Model(&p).Select("id, image_url").Related(&images).Select("name, value").Related(&labels)

	bidStatus := p.GetBidderStatus(userID)

	res := types.Product{
		ID:            p.ID,
		ProductName:   p.ProductName,
		ProductImages: images,
		Desc:          p.Desc,
		Condition:     p.Condition,
		ConditionAvg:  p.ConditionAvg,
		StartPrice:    p.StartPrice,
		BidMultpl:     p.BidMultpl,
		ClosedAT:      p.ClosedAT,
		CreatedAT:     p.CreatedAT,
		Labels:        labels,
		Sold:          p.Sold,
		Closed:        p.Closed,
		BidStatus:     bidStatus,
	}

	return res
}

// ToDetailAPI product detail api type
func (p *Product) ToDetailAPI(userID *int64) types.ProductDetail {
	images := []ProductImage{}
	labels := []ProductLabel{}

	app.DB.Model(&p).Select("image_url").Related(&images).Select("name, value").Related(&labels)

	store := Store{}
	owner := UserSimple{}
	app.DB.First(&store, "id = ?", p.StoreID)
	app.DB.First(&owner, "id = ?", store.OwnerID)
	store.Owner = &owner

	bidStatus := p.GetBidderStatus(userID)

	res := types.ProductDetail{
		ID:            p.ID,
		Store:         store,
		ProductName:   p.ProductName,
		ProductImages: images,
		Desc:          p.Desc,
		Condition:     p.Condition,
		ConditionAvg:  p.ConditionAvg,
		StartPrice:    p.StartPrice,
		BidMultpl:     p.BidMultpl,
		ClosedAT:      p.ClosedAT,
		CreatedAT:     p.CreatedAT,
		Labels:        labels,
		Sold:          p.Sold,
		Closed:        p.Closed,
		BidStatus:     bidStatus,
	}

	return res
}

// ToAPI implementation for ProductBidder
func (p *ProductBidder) ToAPI() *ProductBidder {
	user := UserSimple{}
	app.DB.First(&user, "id = ?", p.UserID)

	p.User = &user

	return p
}
