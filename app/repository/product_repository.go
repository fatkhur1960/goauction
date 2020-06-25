package repository

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/utils"
)

type (
	// ProductRepository init repo
	ProductRepository struct {
		ProductQs models.ProductQuerySet
		ImageQs   models.ProductImageQuerySet
		LabelQs   models.ProductLabelQuerySet
	}

	// NewProductQuery definisi query untuk menambahkan product
	NewProductQuery struct {
		ProductName   string   `json:"product_name" binding:"required"`
		ProductImages []string `json:"product_images"  binding:"required"`
		Desc          string   `json:"desc"  binding:"required"`
		Condition     int32    `json:"condition"  binding:"required"`
		ConditionAvg  float64  `json:"condition_avg" binding:"required"`
		StartPrice    float64  `json:"start_price" binding:"required"`
		BidMultpl     float64  `json:"bid_multpl" binding:"required"`
		ClosedAT      string   `json:"closed_at" binding:"required"`
		Labels        []string `json:"labels" binding:"required"`
	}
)

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

	if err := product.Create(models.DB); err != nil {
		return nil, errors.New("Tidak dapat menambahkan produk")
	}

	for _, url := range query.ProductImages {
		image := models.ProductImage{
			ProductID: product.ID,
			ImageURL:  url,
		}
		image.Create(models.DB)
	}

	for _, name := range query.Labels {
		label := models.ProductLabel{
			ProductID: product.ID,
			Name:      name,
		}
		label.Create(models.DB)
	}

	return product, nil
}

// GetProductList method untuk mendapatkan semua product
func (s *ProductRepository) GetProductList(offset int, limit int, query ...interface{}) (interface{}, int, error) {
	products := []models.Product{}
	conn := s.ProductQs.GetDB()
	count, _ := s.ProductQs.Count()

	// Search product by name
	if len(query) > 0 {
		keyword := fmt.Sprint("%", strings.ToLower(query[0].(string)), "%")
		conn = conn.Where("LOWER(product_name) LIKE ? ", keyword)
		conn = conn.Or("LOWER(\"desc\") LIKE ? ", keyword)

		conn.Count(&count)
	}

	conn.Order("created_at DESC")
	conn.Offset(offset).Limit(limit)
	conn.Find(&products)

	return products, count, nil
}

// GetProductDetail method untuk mendapatkan detail product
func (s *ProductRepository) GetProductDetail(productID int64) (interface{}, error) {
	product := models.Product{}
	err := s.ProductQs.IDEq(productID).GetDB().First(&product)
	if err != nil {
		return nil, err.Error
	}

	return product, nil
}

// CleanUpProduct dao clean all products after testing
// NOTE: using this for testing only
func (s *ProductRepository) CleanUpProduct() {
	products := []models.Product{}
	s.ProductQs.All(&products)
	for _, p := range products {
		if err := p.Delete(models.DB); err != nil {
			log.Fatal(err.Error())
		}
	}
}
