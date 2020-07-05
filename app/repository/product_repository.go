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
		productQs models.ProductQuerySet
		storeQs   models.StoreQuerySet
		imageQs   models.ProductImageQuerySet
		labelQs   models.ProductLabelQuerySet
		bidderQs  models.ProductBidderQuerySet
	}

	// LabelQuery definisi query untuk product label
	LabelQuery struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	// NewProductQuery definisi query untuk menambahkan product
	NewProductQuery struct {
		StoreID       int64        `json:"store_id" binding:"required"`
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
		productQs: models.NewProductQuerySet(app.DB),
		storeQs:   models.NewStoreQuerySet(app.DB),
		imageQs:   models.NewProductImageQuerySet(app.DB),
		labelQs:   models.NewProductLabelQuerySet(app.DB),
		bidderQs:  models.NewProductBidderQuerySet(app.DB),
	}
}

// CreateProduct method untuk menambahkan product
func (s *ProductRepository) CreateProduct(query NewProductQuery) (models.Product, error) {
	closedTime, timeErr := time.Parse(time.RFC3339, query.ClosedAT)
	if timeErr != nil {
		return models.Product{}, errors.New("Invalid datetime format. Correct format is like " + time.RFC3339)
	}

	product := models.Product{}
	product.StoreID = query.StoreID
	product.ProductName = query.ProductName
	product.Desc = query.Desc
	product.Condition = query.Condition
	product.ConditionAvg = query.ConditionAvg
	product.StartPrice = query.StartPrice
	product.BidMultpl = query.BidMultpl
	product.ClosedAT = &closedTime
	product.CreatedAT = &utils.NOW

	if err := product.Create(app.DB); err != nil {
		return models.Product{}, errors.New("Tidak dapat menambahkan produk")
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

	s.storeQs.IDEq(query.StoreID).GetUpdater().SetProductCount(+1).Update()

	return product, nil
}

// GetProductList method untuk mendapatkan semua product
func (s *ProductRepository) GetProductList(args ProductFilter) ([]models.Product, int, error) {
	products := []models.Product{}
	conn := s.productQs.GetDB()
	count, _ := s.productQs.Count()

	conn = conn.Where("sold = ? AND closed = ?", args.Sold, args.Closed)
	// Search product by name
	if args.Query != "" {
		keyword := fmt.Sprint("%", strings.ToLower(args.Query.(string)), "%")
		conn = conn.Where("LOWER(product_name) LIKE ? ", keyword)
		conn = conn.Or("LOWER(\"desc\") LIKE ? ", keyword)

		conn = conn.Count(&count)
	}

	// Search product by user id
	if args.UserID != 0 {
		conn = conn.Where("user_id = ?", args.UserID)
	}

	conn.Order("created_at DESC").Limit(args.Limit).Offset(args.Offset).Find(&products)

	return products, count, nil
}

// GetBidProductList method untuk mendapatkan history bid dari user
func (s *ProductRepository) GetBidProductList(userID int64, offset int, limit int) ([]models.Product, int, error) {
	products := []models.Product{}
	count := 0
	bid := s.bidderQs.GetDB().Select("product_id").Where("user_id = ?", userID).Group("product_id")
	bid.Count(&count)
	err := s.productQs.GetDB().Where("id IN ?", bid.SubQuery()).Offset(offset).Limit(limit).Order("id DESC").Find(&products)
	if err != nil {
		return products, count, err.Error
	}

	return products, count, nil
}

// GetMyProductList method untuk mendapatkan semua product
func (s *ProductRepository) GetMyProductList(args ProductFilter) ([]models.Product, int, error) {
	products := []models.Product{}
	conn := s.productQs.GetDB()
	count, _ := s.productQs.Count()

	conn.Where("user_id = ? AND sold = ? AND closed = ?", args.UserID, args.Sold, args.Closed)
	// Search product by name
	if args.Query != nil {
		keyword := fmt.Sprint("%", strings.ToLower(args.Query.(string)), "%")
		conn = conn.Where("LOWER(product_name) LIKE ? ", keyword)
		conn = conn.Or("LOWER(\"desc\") LIKE ? ", keyword)

		conn.Count(&count)
	}

	conn.Order("created_at DESC").Limit(args.Limit).Offset(args.Offset).Find(&products)

	return products, count, nil
}

// GetByID method untuk mendapatkan product berdasarkan id-nya
func (s *ProductRepository) GetByID(productID int64) (models.Product, error) {
	product := models.Product{}
	err := s.productQs.IDEq(productID).One(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}

// UpdateProduct method untuk mengupdate product
func (s *ProductRepository) UpdateProduct(productID int64, query UpdateProductQuery) (models.Product, error) {
	product := models.Product{}
	dao := s.productQs.IDEq(productID)
	updater := dao.GetUpdater()

	closedTime, _ := time.Parse(time.RFC3339, query.ClosedAT)

	updater.SetProductName(query.ProductName)
	updater.SetDesc(query.Desc)
	updater.SetStartPrice(query.StartPrice)
	updater.SetBidMultpl(query.BidMultpl)
	updater.SetCondition(query.Condition)
	updater.SetConditionAvg(query.ConditionAvg)
	updater.SetClosedAT(&closedTime)
	updater.Update()

	s.imageQs.ProductIDEq(productID).Delete()
	for _, url := range query.ProductImages {
		image := models.ProductImage{
			ProductID: productID,
			ImageURL:  url,
		}
		image.Create(app.DB)
	}

	s.labelQs.ProductIDEq(productID).Delete()
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
func (s *ProductRepository) DeleteProduct(productID int64, storeID int64) error {
	s.labelQs.ProductIDEq(productID).Delete()
	s.imageQs.ProductIDEq(productID).Delete()
	s.storeQs.IDEq(storeID).GetUpdater().SetProductCount(-1).Update()
	return s.productQs.IDEq(productID).Delete()
}

// ReOpenBid digunakan untuk membuka bid kembali
func (s *ProductRepository) ReOpenBid(productID int64, closedAt string) (models.Product, error) {
	product := models.Product{}
	closeTime, timeErr := time.Parse(time.RFC3339, closedAt)
	if timeErr != nil {
		return product, errors.New("Invalid datetime format. Correct format is like " + time.RFC3339)
	}
	dao := s.productQs.IDEq(productID)
	err := dao.GetUpdater().SetClosedAT(&closeTime).SetClosed(false).Update()
	if err != nil {
		return product, err
	}

	dao.One(&product)

	return product, nil
}

// CloseProduct digunakan untuk menutup lelang produk
func (s *ProductRepository) CloseProduct(productID int64) error {
	err := s.productQs.IDEq(productID).GetUpdater().SetClosed(true).Update()
	if err != nil {
		return err
	}

	return nil
}

// SetProductSold digunakan untuk menandai produk sudah terjual
func (s *ProductRepository) SetProductSold(productID int64) (models.Product, error) {
	product := models.Product{}
	dao := s.productQs.IDEq(productID)
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
	s.productQs.All(&products)
	for _, p := range products {
		if err := p.Delete(app.DB); err != nil {
			log.Fatal(err.Error())
		}
	}
}
