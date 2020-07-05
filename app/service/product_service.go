package service

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	mid "github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/models"
	repo "github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/app/types"
	"github.com/fatkhur1960/goauction/system/event"
	"github.com/fatkhur1960/goauction/system/queue"
	"github.com/gin-gonic/gin"
)

type (
	// ProductService api product implementation
	ProductService struct {
		sync.Mutex
		productRepo *repo.ProductRepository
		storeRepo   *repo.StoreRepository
		event       *event.Listener
	}

	// BidProductQuery query untuk bid product
	BidProductQuery struct {
		ProductID int64   `json:"product_id" binding:"required"`
		BidPrice  float64 `json:"bid_price" binding:"required"`
	}

	// ReOpenBidQuery query untuk membuka bid lagi
	ReOpenBidQuery struct {
		ProductID int64  `json:"product_id" binding:"required"`
		ClosedAT  string `json:"closed_at" binding:"required"`
	}
)

// NewProductService api instance
// @RouterGroup /product/v1
func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repo.NewProductRepository(),
		storeRepo:   repo.NewStoreRepository(),
		event:       event.NewListener(queue.JobQueue),
	}
}

// AddProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk menambahkan product
// @Accept json
// @Param product_name body string true "ProductName"
// @Param product_images body []string true "ProductImages"
// @Param desc body string true "Desc"
// @Param condition body int true "Condition"
// @Param condition_avg body int true "ConditionAvg"
// @Param start_price body int true "StartPrice"
// @Param bid_multpl body int true "BidMultpl"
// @Param closed_at body string true "ClosedAt"
// @Param labels body []string true "Labels"
// @Produce json
// @Success 200 {object} app.Result{result=models.Product}
// @Failure 400 {object} app.Result
// @Router /add [post] [auth]
func (s *ProductService) AddProduct(c *gin.Context, query *repo.NewProductQuery) {
	store, e := s.storeRepo.GetByID(query.StoreID)
	if e != nil {
		APIResult.Error(c, http.StatusBadRequest, "Store ID tidak valid")
		return
	} else if store.OwnerID != mid.CurrentUser.ID {
		APIResult.Error(c, http.StatusBadRequest, "Anda tidak dapat menambahkan product ke store ini")
		return
	}

	product, err := s.productRepo.CreateProduct(*query)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, product.ToAPI(&mid.CurrentUser.ID))
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
// @Param filter query string false "Filter"
// @Success 200 {object} app.Result{result=EntriesResult{entries=[]types.Product}}
// @Failure 400 {object} app.Result
// @Router /list [get]
func (s *ProductService) ListProduct(c *gin.Context, query *QueryEntries) {
	userID := int64(0)
	closed := false
	sold := false
	if query.Filter != "" {
		args := strings.Split(query.Filter, ",")
		for _, arg := range args {
			if strings.HasPrefix(arg, "user_id:") {
				res, _ := strconv.Atoi(arg[len("user_id:"):])
				userID = int64(res)
			} else if strings.HasPrefix(arg, "closed:") {
				res, _ := strconv.ParseBool(arg[len("closed:"):])
				closed = res
			} else if strings.HasPrefix(arg, "sold:") {
				res, _ := strconv.ParseBool(arg[len("sold:"):])
				sold = res
			}
		}
	}

	filter := repo.ProductFilter{
		Query:  query.Query,
		Offset: query.Offset,
		Limit:  query.Limit,
		UserID: userID,
		Sold:   sold,
		Closed: closed,
	}

	products, count, _ := s.productRepo.GetProductList(filter)
	entries := []types.Product{}
	for _, product := range products {
		entries = append(entries, product.ToAPI(&mid.CurrentUser.ID))
	}

	APIResult.Success(c, EntriesResult{entries, count})
}

// ListMyProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk mendapatkan list product untuk current user
// @Produce json
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Param query query string false "Query"
// @Param filter query string false "Filter"
// @Success 200 {object} app.Result{result=EntriesResult{entries=[]types.Product}}
// @Failure 400 {object} app.Result
// @Router /me/list [get] [auth]
func (s *ProductService) ListMyProduct(c *gin.Context, query *QueryEntries) {
	entries := []types.Product{}
	closed, sold := false, false

	if query.Filter != "" {
		if strings.Contains(query.Filter, "sold") {
			sold = true
		} else if strings.Contains(query.Filter, "closed") {
			closed = true
		}
	}

	filter := repo.ProductFilter{
		UserID: mid.CurrentUser.ID,
		Closed: closed,
		Sold:   sold,
		Query:  query.Query,
		Offset: query.Offset,
		Limit:  query.Limit,
	}
	products, count, _ := s.productRepo.GetMyProductList(filter)

	for _, product := range products {
		entries = append(entries, product.ToAPI(&mid.CurrentUser.ID))
	}

	APIResult.Success(c, EntriesResult{entries, count})
}

// DetailProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk menampilkan detail product
// @Accept json
// @Produce json
// @Param id query int true "ID"
// @Success 200 {object} app.Result{result=types.ProductDetail}
// @Failure 400 {object} app.Result
// @Router /detail [get] [auth]
func (s *ProductService) DetailProduct(c *gin.Context, query *IDQuery) {
	product, err := s.productRepo.GetByID(query.ID)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, product.ToDetailAPI())
}

// UpdateProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk mengupdate product
// @Accept json
// @Produce json
// @Param id body int true "ID"
// @Param product_name body string true "ProductName"
// @Param product_images body []string true "ProductImages"
// @Param desc body string true "Desc"
// @Param condition body int true "Condition"
// @Param condition_avg body int true "ConditionAvg"
// @Param start_price body int true "StartPrice"
// @Param bid_multpl body int true "BidMultpl"
// @Param closed_at body string true "ClosedAt"
// @Param labels body []string true "Labels"
// @Success 200 {object} app.Result{result=models.Product}
// @Failure 400 {object} app.Result
// @Router /update [post] [auth]
func (s *ProductService) UpdateProduct(c *gin.Context, query *repo.UpdateProductQuery) {
	p, err := s.productRepo.GetByID(query.ID)
	store, _ := s.storeRepo.GetByID(p.StoreID)

	updateTime, parseTimeError := time.Parse(time.RFC3339, query.ClosedAT)

	if err != nil {
		APIResult.Error(c, http.StatusNoContent, "Produk tidak ditemukan")
		return
	} else if store.OwnerID != mid.CurrentUser.ID {
		APIResult.Error(c, http.StatusBadRequest, "Unauthorized")
		return
	} else if p.Closed {
		APIResult.Error(c, http.StatusBadRequest, "Bid sudah ditutup, buka lagi untuk bisa mengupdate")
		return
	} else if updateTime.Sub(*p.ClosedAT) < 0 {
		APIResult.Error(c, http.StatusBadRequest, "Waktu ditutup tidak valid")
		return
	}

	product, err := s.productRepo.UpdateProduct(query.ID, *query)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	} else if parseTimeError != nil {
		log.Fatalf("ParseTime] error parsing time %v", parseTimeError.Error())
	}

	APIResult.Success(c, product.ToAPI(&mid.CurrentUser.ID))
}

// DeleteProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk menghapus product
// @Accept json
// @Produce json
// @Param id body int true "ID"
// @Success 200 {object} app.Result{result=models.Product}
// @Failure 400 {object} app.Result
// @Router /delete [post] [auth]
func (s *ProductService) DeleteProduct(c *gin.Context, query *IDQuery) {
	product, e := s.productRepo.GetByID(query.ID)
	store, _ := s.storeRepo.GetByID(product.StoreID)

	if store.OwnerID != mid.CurrentUser.ID {
		APIResult.Error(c, http.StatusBadRequest, "Anda tidak dapat menghapus produk ini")
		return
	} else if e != nil {
		APIResult.Error(c, http.StatusBadRequest, "Produk tidak ditemukan")
		return
	}

	if err := s.productRepo.DeleteProduct(product.ID, store.ID); err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, nil)
}

// BidProduct docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk mengupdate product
// @Accept json
// @Produce json
// @Param product_id body int true "ProductID"
// @Param bid_price body number true "BidPrice"
// @Success 200 {object} app.Result{result=models.ProductBidder}
// @Failure 400 {object} app.Result
// @Router /bid [post] [auth]
func (s *ProductService) BidProduct(c *gin.Context, query *BidProductQuery) {
	product, err1 := s.productRepo.GetByID(query.ProductID)
	store, _ := s.storeRepo.GetByID(product.StoreID)

	if store.OwnerID == mid.CurrentUser.ID {
		APIResult.Error(c, http.StatusBadRequest, "Anda tidak dapat melakukan bid ini")
		return
	} else if product.Closed {
		APIResult.Error(c, http.StatusBadRequest, "Bid sudah ditutup")
		return
	} else if err1 != nil {
		APIResult.Error(c, http.StatusBadRequest, "Bid tidak ditemukan")
		return
	} else if (int(query.BidPrice) % int(product.BidMultpl)) != 0 {
		APIResult.Error(c, http.StatusBadRequest, fmt.Sprintf("Bid tidak termasuk kelipatan %v", product.BidMultpl))
		return
	} else if query.BidPrice <= product.GetLatestBidPrice() {
		APIResult.Error(c, http.StatusBadRequest, fmt.Sprintf("Bid harus lebih besar dari %v", product.GetLatestBidPrice()))
		return
	}

	bidder, err := s.productRepo.AddProductBidder(mid.CurrentUser.ID, query.ProductID, query.BidPrice)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, fmt.Sprintf("Error: %s", err.Error()))
	}

	{
		go s.event.Emmit(&event.UserBidProductEvent{
			User:    &mid.CurrentUser,
			Product: product,
			BidData: bidder.(models.ProductBidder),
		})
	}

	APIResult.Success(c, bidder)
}

// ReOpenProductBid docs
// @Tags ProductService
// @Summary Endpoint digunakan untuk membuka bid kembali
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id body int true "ProductID"
// @Param closed_at body string true "ClosedAT"
// @Success 200 {object} app.Result{result=models.Product}
// @Failure 400 {object} app.Result
// @Router /reopen [post] [auth]
func (s *ProductService) ReOpenProductBid(c *gin.Context, query *ReOpenBidQuery) {
	p, err := s.productRepo.GetByID(query.ProductID)
	store, _ := s.storeRepo.GetByID(p.StoreID)
	updatedTime, parseTimeError := time.Parse(time.RFC3339, query.ClosedAT)

	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, "Produk tidak ditemukan")
		return
	} else if store.OwnerID != mid.CurrentUser.ID {
		APIResult.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	} else if !p.Closed {
		APIResult.Error(c, http.StatusBadRequest, "Bid belum ditutup")
		return
	} else if updatedTime.Sub(*p.ClosedAT) <= 0 {
		APIResult.Error(c, http.StatusBadRequest, "Waktu ditutup tidak valid")
		return
	}

	product, err2 := s.productRepo.ReOpenBid(query.ProductID, query.ClosedAT)
	if err2 != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	} else if parseTimeError != nil {
		log.Fatalf("ParseTime] error parsing time %v", parseTimeError.Error())
	}

	APIResult.Success(c, product)
}

// MarkProductAsSold docs
// @Tags ProductService
// @Summary Endpoint digunakan untuk menandai produk sudah terjual
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id body int true "ID"
// @Success 200 {object} app.Result{result=models.Product}
// @Failure 400 {object} app.Result
// @Router /mark-as-sold [post] [auth]
func (s *ProductService) MarkProductAsSold(c *gin.Context, query *IDQuery) {
	p, err := s.productRepo.GetByID(query.ID)
	store, _ := s.storeRepo.GetByID(p.StoreID)

	if err != nil {
		APIResult.Error(c, http.StatusNoContent, "Produk tidak ditemukan")
		return
	} else if store.OwnerID != mid.CurrentUser.ID {
		APIResult.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	} else if !p.Closed {
		APIResult.Error(c, http.StatusBadRequest, "Bid belum ditutup")
		return
	} else if p.Sold {
		APIResult.Error(c, http.StatusBadRequest, "Sudah terjual")
		return
	}

	product, err := s.productRepo.SetProductSold(query.ID)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, product)
}
