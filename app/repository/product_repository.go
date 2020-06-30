package repository

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/utils"
)

type (
	// ProductRepository init repo
	ProductRepository struct {
		ProductQs models.ProductQuerySet
		ImageQs   models.ProductImageQuerySet
		LabelQs   models.ProductLabelQuerySet
		BidderQs  models.ProductBidderQuerySet
	}

	// LabelQuery definisi query untuk product label
	LabelQuery struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	// NewProductQuery definisi query untuk menambahkan product
	NewProductQuery struct {
		ProductName   string       `json:"product_name" binding:"required"`
		ProductImages []string     `json:"product_images"  binding:"required"`
		Desc          string       `json:"desc"  binding:"required"`
		Condition     int32        `json:"condition"  binding:"required"`
		ConditionAvg  float64      `json:"condition_avg" binding:"required"`
		StartPrice    float64      `json:"start_price" binding:"required"`
		BidMultpl     float64      `json:"bid_multpl" binding:"required"`
		ClosedAT      string       `json:"closed_at" binding:"required"`
		Labels        []LabelQuery `json:"labels" binding:"required"`
	}

	// UpdateProductQuery definisi query untuk menambahkan product
	UpdateProductQuery struct {
		ID            int64        `json:"id" binding:"required"`
		ProductName   string       `json:"product_name" binding:"required"`
		ProductImages []string     `json:"product_images"  binding:"required"`
		Desc          string       `json:"desc"  binding:"required"`
		Condition     int32        `json:"condition"  binding:"required"`
		ConditionAvg  float64      `json:"condition_avg" binding:"required"`
		StartPrice    float64      `json:"start_price" binding:"required"`
		BidMultpl     float64      `json:"bid_multpl" binding:"required"`
		ClosedAT      string       `json:"closed_at" binding:"required"`
		Labels        []LabelQuery `json:"labels" binding:"required"`
	}

	// ProductFilter definisi type untuk filter product
	ProductFilter struct {
		UserID int64
		Sold   bool
		Closed bool
		Query  interface{}
		Offset int
		Limit  int
	}
)

// NewProductRepository create instance for product repository
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		ProductQs: models.NewProductQuerySet(app.DB),
		ImageQs:   models.NewProductImageQuerySet(app.DB),
		LabelQs:   models.NewProductLabelQuerySet(app.DB),
		BidderQs:  models.NewProductBidderQuerySet(app.DB),
	}
}

// CreateProduct method untuk menambahkan product
func (s *ProductRepository) CreateProduct(userID int64, query NewProductQuery) (interface{}, error) {
	closedTime, timeErr := time.Parse(time.RFC3339, query.ClosedAT)
	if timeErr != nil {
		return nil, errors.New("Invalid datetime format. Correct format is like " + time.RFC3339)
	}

	product := models.Product{}
	product.UserID = userID
	product.ProductName = query.ProductName
	product.Desc = query.Desc
	product.Condition = query.Condition
	product.ConditionAvg = query.ConditionAvg
	product.StartPrice = query.StartPrice
	product.BidMultpl = query.BidMultpl
	product.ClosedAT = &closedTime
	product.CreatedAT = &utils.NOW

	if err := product.Create(app.DB); err != nil {
		return nil, errors.New("Tidak dapat menambahkan produk")
	}

	for _, url := range query.ProductImages {
		image := models.ProductImage{
			ProductID: product.ID,
			ImageURL:  url,
		}
		image.Create(app.DB)
	}

	for _, label := range query.Labels {
		label := models.ProductLabel{
			ProductID: product.ID,
			Name:      label.Name,
			Value:     label.Value,
		}
		label.Create(app.DB)
	}

	return product, nil
}

// GetProductList method untuk mendapatkan semua product
func (s *ProductRepository) GetProductList(args ProductFilter) ([]models.Product, int, error) {
	products := []models.Product{}
	conn := s.ProductQs.GetDB()
	count, _ := s.ProductQs.Count()

	conn.Where("sold = ? AND closed = ?", args.Sold, args.Closed)
	// Search product by name
	if args.Query != nil {
		keyword := fmt.Sprint("%", strings.ToLower(args.Query.(string)), "%")
		conn = conn.Where("LOWER(product_name) LIKE ? ", keyword)
		conn = conn.Or("LOWER(\"desc\") LIKE ? ", keyword)

		conn.Count(&count)
	}

	conn.Order("created_at DESC")
	conn.Offset(args.Offset).Limit(args.Limit)
	conn.Find(&products)

	return products, count, nil
}

// GetBidProductList method untuk mendapatkan history bid dari user
func (s *ProductRepository) GetBidProductList(userID int64, offset int, limit int) ([]models.Product, int, error) {
	products := []models.Product{}
	count := 0
	bid := s.BidderQs.GetDB().Select("product_id").Where("user_id = ?", userID).Group("product_id")
	bid.Count(&count)
	err := s.ProductQs.GetDB().Where("id IN ?", bid.SubQuery()).Offset(offset).Limit(limit).Order("id DESC").Find(&products)
	if err != nil {
		return products, count, err.Error
	}

	return products, count, nil
}

// GetMyProductList method untuk mendapatkan semua product
func (s *ProductRepository) GetMyProductList(args ProductFilter) ([]models.Product, int, error) {
	products := []models.Product{}
	conn := s.ProductQs.GetDB()
	count, _ := s.ProductQs.Count()

	conn.Where("user_id = ? AND sold = ? AND closed = ?", args.UserID, args.Sold, args.Closed)
	// Search product by name
	if args.Query != nil {
		keyword := fmt.Sprint("%", strings.ToLower(args.Query.(string)), "%")
		conn = conn.Where("LOWER(product_name) LIKE ? ", keyword)
		conn = conn.Or("LOWER(\"desc\") LIKE ? ", keyword)

		conn.Count(&count)
	}

	conn.Order("created_at DESC")
	conn.Offset(args.Offset).Limit(args.Limit)
	conn.Find(&products)

	return products, count, nil
}

// GetByID method untuk mendapatkan product berdasarkan id-nya
func (s *ProductRepository) GetByID(productID int64) (models.Product, error) {
	product := models.Product{}
	err := s.ProductQs.IDEq(productID).One(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}

// UpdateProduct method untuk mengupdate product
func (s *ProductRepository) UpdateProduct(productID int64, query UpdateProductQuery) (models.Product, error) {
	product := models.Product{}
	dao := s.ProductQs.IDEq(productID)
	updater := dao.GetUpdater()

	closedTime, _ := time.Parse(time.RFC3339, query.ClosedAT)

	updater.SetProductName(query.ProductName)
	updater.SetDesc(query.Desc)
	updater.SetCondition(query.Condition)
	updater.SetConditionAvg(query.ConditionAvg)
	updater.SetClosedAT(&closedTime)
	updater.Update()

	s.ImageQs.ProductIDEq(productID).Delete()
	for _, url := range query.ProductImages {
		image := models.ProductImage{
			ProductID: productID,
			ImageURL:  url,
		}
		image.Create(app.DB)
	}

	s.LabelQs.ProductIDEq(productID).Delete()
	for _, label := range query.Labels {
		label := models.ProductLabel{
			ProductID: productID,
			Name:      label.Name,
			Value:     label.Value,
		}
		label.Create(app.DB)
	}

	dao.One(&product)

	return product, nil
}

// AddProductBidder digunakan untuk menyimpan user bid product
func (s *ProductRepository) AddProductBidder(userID int64, productID int64, bidPrice float64) (interface{}, error) {
	bidder := models.ProductBidder{
		UserID:    userID,
		ProductID: productID,
		BidPrice:  bidPrice,
	}

	if err := bidder.Create(app.DB); err != nil {
		return nil, err
	}

	return bidder, nil
}

// DeleteProduct digunakan untuk menghapus product
func (s *ProductRepository) DeleteProduct(productID int64) error {
	s.LabelQs.ProductIDEq(productID).Delete()
	s.ImageQs.ProductIDEq(productID).Delete()
	return s.ProductQs.IDEq(productID).Delete()
}

// ReOpenBid digunakan untuk membuka bid kembali
func (s *ProductRepository) ReOpenBid(productID int64, closedAt string) (models.Product, error) {
	product := models.Product{}
	closeTime, timeErr := time.Parse(time.RFC3339, closedAt)
	if timeErr != nil {
		return product, errors.New("Invalid datetime format. Correct format is like " + time.RFC3339)
	}
	dao := s.ProductQs.IDEq(productID)
	err := dao.GetUpdater().SetClosedAT(&closeTime).SetClosed(false).Update()
	if err != nil {
		return product, err
	}

	dao.One(&product)

	return product, nil
}

// SetProductSold digunakan untuk menandai produk sudah terjual
func (s *ProductRepository) SetProductSold(productID int64) (models.Product, error) {
	product := models.Product{}
	dao := s.ProductQs.IDEq(productID)
	err := dao.GetUpdater().SetSold(true).Update()
	if err != nil {
		return product, err
	}

	dao.One(&product)

	return product, nil
}

// CleanUpProduct dao clean all products after testing
// NOTE: using this for testing only
func (s *ProductRepository) CleanUpProduct() {
	products := []models.Product{}
	s.ProductQs.All(&products)
	for _, p := range products {
		if err := p.Delete(app.DB); err != nil {
			log.Fatal(err.Error())
		}
	}
}
