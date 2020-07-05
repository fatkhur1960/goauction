package test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/app/service"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/fatkhur1960/goauction/tests/endpoint"
	"github.com/go-playground/assert/v2"
	"github.com/mitchellh/mapstructure"
	"syreclabs.com/go/faker"
)

func TestAddProduct(t *testing.T) {
	// defer cleanUsers()
	labels := []repository.LabelQuery{}
	labels = append(labels, repository.LabelQuery{
		Name:  "label_name",
		Value: "value",
	})
	token := authorizeUser()
	store := upgradeUser(token)
	payload := repository.NewProductQuery{
		StoreID:       store.ID,
		ProductName:   faker.Commerce().ProductName(),
		ProductImages: []string{faker.Internet().Url()},
		Desc:          faker.RandomString(100),
		Condition:     1,
		ConditionAvg:  100,
		StartPrice:    float64(faker.Commerce().Price()),
		BidMultpl:     float64(faker.Commerce().Price()),
		ClosedAT:      utils.NOW.Add(time.Hour * 24).Format(time.RFC3339),
		Labels:        labels,
	}

	rv := reqPOST(endpoint.AddProduct, payload, token)
	assert.Equal(t, rv.Code, 0)
	resMap := rv.Result.(map[string]interface{})
	assert.NotEqual(t, resMap["product_name"], nil)
	assert.NotEqual(t, resMap["product_images"], nil)
	assert.NotEqual(t, resMap["desc"], nil)
	assert.NotEqual(t, resMap["condition"], nil)
	assert.NotEqual(t, resMap["condition_avg"], nil)
	assert.NotEqual(t, resMap["start_price"], nil)
	assert.NotEqual(t, resMap["bid_multpl"], nil)
	assert.NotEqual(t, resMap["closed_at"], nil)
	assert.NotEqual(t, resMap["labels"], nil)
	assert.NotEqual(t, resMap["sold"], true)
}

func TestListProduct(t *testing.T) {
	// defer cleanUsers()
	token := authorizeUser()
	store := upgradeUser(token)
	createProduct(token, store.ID)
	createProduct(token, store.ID)
	_, err := createProduct(token, store.ID)

	if err != nil {
		t.SkipNow()
	}

	res := service.EntriesResult{}

	rv := reqGET(endpoint.ListProduct+"?limit=10&offset=0", token)
	assert.Equal(t, rv.Code, 0)
	resMap := rv.Result.(map[string]interface{})
	mapstructure.Decode(resMap, &res)
	assert.Equal(t, res.Count, 4)
}

func TestBidProduct(t *testing.T) {
	token := authorizeUser()
	store := upgradeUser(token)
	product, _ := createProduct(token, store.ID)
	token2 := authorizeUser()
	payload := service.BidProductQuery{
		ProductID: product.ID,
		BidPrice:  50000,
	}
	rv := reqPOST(endpoint.BidProduct, payload, token2)
	assert.Equal(t, rv.Code, 0)
}

func TestBidProductWithSameUser(t *testing.T) {
	token := authorizeUser()
	store := upgradeUser(token)
	product, _ := createProduct(token, store.ID)
	payload := service.BidProductQuery{
		ProductID: product.ID,
		BidPrice:  50000,
	}
	rv := reqPOST(endpoint.BidProduct, payload, token)
	assert.Equal(t, rv.Description, "Anda tidak dapat melakukan bid ini")
}

func TestBidProductClosed(t *testing.T) {
	token := authorizeUser()
	store := upgradeUser(token)
	product, _ := createProduct(token, store.ID)
	closeProduct(product.ID)
	token2 := authorizeUser()
	payload := service.BidProductQuery{
		ProductID: product.ID,
		BidPrice:  50000,
	}
	rv := reqPOST(endpoint.BidProduct, payload, token2)
	assert.Equal(t, rv.Description, "Bid sudah ditutup")
}

func TestBidWithSamePrice(t *testing.T) {
	token := authorizeUser()
	store := upgradeUser(token)
	product, _ := createProduct(token, store.ID)
	token2 := authorizeUser()
	payload := service.BidProductQuery{
		ProductID: product.ID,
		BidPrice:  50000,
	}
	rv := reqPOST(endpoint.BidProduct, payload, token2)
	assert.Equal(t, rv.Code, 0)
	rv2 := reqPOST(endpoint.BidProduct, payload, token2)
	assert.Equal(t, rv2.Description, fmt.Sprintf("Bid harus lebih besar dari %v", payload.BidPrice))
}

func TestBidWithInvalidPrice(t *testing.T) {
	token := authorizeUser()
	store := upgradeUser(token)
	product, _ := createProduct(token, store.ID)
	token2 := authorizeUser()
	payload := service.BidProductQuery{
		ProductID: product.ID,
		BidPrice:  60000,
	}
	rv := reqPOST(endpoint.BidProduct, payload, token2)
	assert.Equal(t, rv.Description, fmt.Sprintf("Bid tidak termasuk kelipatan %v", product.BidMultpl))
}

func TestBidWithProductNotFound(t *testing.T) {
	token := authorizeUser()
	payload := service.BidProductQuery{
		ProductID: 99999,
		BidPrice:  1000000,
	}
	rv := reqPOST(endpoint.BidProduct, payload, token)
	assert.Equal(t, rv.Description, "Bid tidak ditemukan")
}

func TestUpdateProduct(t *testing.T) {
	token := authorizeUser()
	store := upgradeUser(token)
	product, _ := createProduct(token, store.ID)
	labels := []repository.LabelQuery{}
	labels = append(labels, repository.LabelQuery{
		Name:  "label_name",
		Value: "value",
	})
	payload := repository.UpdateProductQuery{
		ID:            product.ID,
		ProductName:   faker.Commerce().ProductName(),
		ProductImages: []string{faker.Internet().Url()},
		Desc:          faker.RandomString(100),
		Condition:     2,
		ConditionAvg:  90,
		StartPrice:    50000.0,
		BidMultpl:     50000.0,
		ClosedAT:      utils.NOW.Add(time.Hour * 24).Format(time.RFC3339),
		Labels:        labels,
	}

	rv := reqPOST(endpoint.UpdateProduct, payload, token)
	assert.Equal(t, rv.Code, 0)
	resMap := rv.Result.(map[string]interface{})
	assert.Equal(t, resMap["product_name"], payload.ProductName)
	// assert.Equal(t, resMap["product_images"], payload.ProductImages)
	assert.Equal(t, resMap["desc"], payload.Desc)
	assert.Equal(t, resMap["condition"], float64(payload.Condition))
	assert.Equal(t, resMap["condition_avg"], payload.ConditionAvg)
	assert.Equal(t, resMap["start_price"], payload.StartPrice)
	assert.Equal(t, resMap["bid_multpl"], payload.BidMultpl)
	assert.Equal(t, resMap["closed_at"], payload.ClosedAT)
	labelLeft, _ := json.Marshal(resMap["labels"])
	labelRight, _ := json.Marshal(payload.Labels)
	assert.Equal(t, labelLeft, labelRight)
}

func TestDeleteProduct(t *testing.T) {
	token := authorizeUser()
	store := upgradeUser(token)
	product, _ := createProduct(token, store.ID)
	id := strconv.Itoa(int(product.ID))
	payload := service.IDQuery{
		ID: product.ID,
	}
	rv := reqPOST(endpoint.DeleteProduct, payload, token)
	assert.Equal(t, rv.Code, 0)

	rv2 := reqGET(endpoint.DetailProduct+"?id="+id, token)
	assert.Equal(t, rv2.Code, 4000)
}

func TestReOpenProduct(t *testing.T) {
	token := authorizeUser()
	store := upgradeUser(token)
	product, _ := createProduct(token, store.ID)
	closeProduct(product.ID)
	payload := service.ReOpenBidQuery{
		ProductID: product.ID,
		ClosedAT:  utils.NOW.Add(time.Hour * 24 * 2).Format(time.RFC3339),
	}
	rv := reqPOST(endpoint.ReOpenProductBid, payload, token)
	assert.Equal(t, rv.Code, 0)
}

func TestReOpenProductNotFound(t *testing.T) {
	token := authorizeUser()
	payload := service.ReOpenBidQuery{
		ProductID: 1001,
		ClosedAT:  utils.NOW.Add(time.Hour * 24 * 2).Format(time.RFC3339),
	}
	rv := reqPOST(endpoint.ReOpenProductBid, payload, token)
	assert.Equal(t, rv.Description, "Produk tidak ditemukan")
}

func TestReOpenProductWithOtherUser(t *testing.T) {
	token := authorizeUser()
	store := upgradeUser(token)
	product, _ := createProduct(token, store.ID)
	token2 := authorizeUser()
	payload := service.ReOpenBidQuery{
		ProductID: product.ID,
		ClosedAT:  utils.NOW.Add(time.Hour * 24 * 2).Format(time.RFC3339),
	}
	rv := reqPOST(endpoint.ReOpenProductBid, payload, token2)
	assert.Equal(t, rv.Description, "Unauthorized")
}

func TestReOpenProductWithInvalidCloseTime(t *testing.T) {
	token := authorizeUser()
	store := upgradeUser(token)
	product, _ := createProduct(token, store.ID)
	closeProduct(product.ID)
	payload := service.ReOpenBidQuery{
		ProductID: product.ID,
		ClosedAT:  utils.NOW.Add(-(time.Hour * 24)).Format(time.RFC3339),
	}
	rv := reqPOST(endpoint.ReOpenProductBid, payload, token)
	assert.Equal(t, rv.Description, "Waktu ditutup tidak valid")
}
