package test

import (
	"strings"
	"testing"
	"time"

	"github.com/fatkhur1960/goauction/app/models"
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
	token := authorizeUser()
	payload := repository.NewProductQuery{
		ProductName:   faker.Commerce().ProductName(),
		ProductImages: []string{faker.Internet().Url()},
		Desc:          faker.RandomString(100),
		Condition:     1,
		ConditionAvg:  100,
		StartPrice:    float64(faker.Commerce().Price()),
		BidMultpl:     float64(faker.Commerce().Price()),
		ClosedAT:      utils.NOW.Add(time.Hour * 24).Format(time.RFC3339),
		Labels:        []repository.LabelQuery{},
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
	createProduct(token)
	createProduct(token)
	_, err := createProduct(token)

	if err != nil {
		t.SkipNow()
	}

	res := service.EntriesResult{}

	rv := reqGET(endpoint.ListProduct+"?limit=0&offset=10", token)
	assert.Equal(t, rv.Code, 0)
	resMap := rv.Result.(map[string]interface{})
	mapstructure.Decode(resMap, &res)
	assert.Equal(t, res.Count, 4)
}

func TestBidProduct(t *testing.T) {
	assert.Equal(t, true, true)
}

func TestUpdateProduct(t *testing.T) {
	token := authorizeUser()
	p, _ := createProduct(token)
	product := p.(*models.Product)

	path := strings.ReplaceAll(endpoint.UpdateProduct, ":id", string(product.ID))
	rv := reqPOST(path, token)
	assert.Equal(t, rv.Code, 0)
	// resMap := rv.Result.(map[string]interface{})
	// t.Log(resMap)
}

func TestDeleteProduct(t *testing.T) {
	assert.Equal(t, true, true)
}
