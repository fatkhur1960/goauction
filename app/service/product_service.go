package service

import (
	"net/http"

	mid "github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/models"
	repo "github.com/fatkhur1960/goauction/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type (
	// ProductService api product implementation
	ProductService struct {
		productRepo repo.ProductRepository
	}
)

// NewProductService api instance
// @RouterGroup /product/v1
func NewProductService(db *gorm.DB) ProductService {
	return ProductService{
		productRepo: repo.ProductRepository{
			ProductQs: models.NewProductQuerySet(db),
			ImageQs:   models.NewProductImageQuerySet(db),
			LabelQs:   models.NewProductLabelQuerySet(db),
		},
	}
}

// AddProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk menambahkan product
// @Accept json
// @Produce json
// @Success 200 {object} app.Result{result=models.Product}
// @Failure 400 {object} app.Result
// @Router /add [post] [auth]
func (s *ProductService) AddProduct(c *gin.Context) {
	query := repo.NewProductQuery{}
	validateRequest(c, &query)

	rv, err := s.productRepo.CreateProduct(mid.CurrentUser.ID, query)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product := Product{}
	APIResult.Success(c, product.toAPI(rv, models.DB))
}

// ListProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk menampilkan list product
// @Accept json
// @Produce json
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Param query query string false "Query"
// @Success 200 {object} app.Result{result=EntriesResult{entries=[]models.Product}}
// @Failure 400 {object} app.Result
// @Router /list [get] [auth]
func (s *ProductService) ListProduct(c *gin.Context) {
	query := &QueryEntries{}
	query.validate(c, typeQuery)

	var entries []interface{}
	products, count, _ := s.productRepo.GetProductList(query.Offset, query.Limit, query.Query)

	for _, item := range products.([]models.Product) {
		product := Product{}
		entries = append(entries, product.toAPI(item, models.DB))
	}

	APIResult.Success(c, EntriesResult{entries, count})
}

// DetailProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk menampilkan detail product
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} app.Result{result=EntriesResult{entries=[]models.Product}}
// @Failure 400 {object} app.Result
// @Router /detail/:id [get] [auth]
func (s *ProductService) DetailProduct(c *gin.Context) {
	query := &IDQuery{}
	query.validate(c, typeURI)

	rv, err := s.productRepo.GetProductDetail(query.ID)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product := Product{}

	APIResult.Success(c, product.toAPI(rv, models.DB))
}

// UpdateProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk mengupdate product
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} app.Result{result=models.Product}
// @Failure 400 {object} app.Result
// @Router /update/:id [post] [auth]
func (s *ProductService) UpdateProduct(c *gin.Context) {

}

// DeleteProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk menghapus product
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} app.Result{result=models.Product}
// @Failure 400 {object} app.Result
// @Router /delete/:id [post] [auth]
func (s *ProductService) DeleteProduct(c *gin.Context) {

}

// BidProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk mengupdate product
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} app.Result{result=models.Product}
// @Failure 400 {object} app.Result
// @Router /bid/:id [post] [auth]
func (s *ProductService) BidProduct(c *gin.Context) {

}
