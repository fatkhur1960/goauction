package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/app/router"
	"github.com/fatkhur1960/goauction/app/service"
	"github.com/fatkhur1960/goauction/app/types"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/fatkhur1960/goauction/tests/endpoint"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"syreclabs.com/go/faker"
)

var (
	ts     = httptest.NewServer(getTestingRoutes())
	client = http.Client{}
)

func getTestingRoutes() *gin.Engine {
	app.ConnectDatabaseTest()
	gin.SetMode(gin.TestMode)
	router := router.GetGeneratedRoutes(gin.New())
	return router
}

func parseResult(resp *http.Response, err error) app.Result {
	var result app.Result
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(responseData, &result)

	if os.Getenv("TEST_LOG") == "debug" {
		res, _ := json.MarshalIndent(&result, "", "\t")
		log.Printf("=> result %s\n", string(res))
	}

	return result
}

func reqPOST(path string, args ...interface{}) app.Result {
	url := fmt.Sprintf("%s/api%s", ts.URL, path)
	payload := args[0]
	var token string
	if len(args) > 1 {
		token = args[1].(string)
	}

	query, _ := json.Marshal(payload)
	if os.Getenv("TEST_LOG") == "debug" {
		log.Printf("=> sending payload to %s %s\n", url, string(query))
	}
	req, reqErr := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(query),
	)
	if reqErr != nil {
		log.Println("]", reqErr.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	resp, err := client.Do(req)

	return parseResult(resp, err)
}

func reqGET(args ...interface{}) app.Result {
	path := args[0]
	var token string
	if len(args) > 1 {
		token = args[1].(string)
	}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api%s", ts.URL, path),
		nil,
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)

	return parseResult(resp, err)
}

func generateUserThenActivate() (int64, string, string) {
	u := service.RegisterUserQuery{
		FullName: faker.Name().Name(),
		Email:    faker.Internet().Email(),
		PhoneNum: faker.PhoneNumber().CellPhone(),
	}

	passhash := faker.Internet().Password(8, 12)

	rv := reqPOST(endpoint.RegisterUser, u)
	rMap := rv.Result.(map[string]interface{})

	activate := service.ActivateUserQuery{
		Token:    fmt.Sprintf("%v", rMap["token"]),
		Passhash: passhash,
	}
	rv2 := reqPOST(endpoint.ActivateUser, activate)
	resMap := rv2.Result.(map[string]interface{})
	userID := resMap["user_id"].(float64)

	return int64(userID), u.Email, passhash
}

func authorizeUser() string {
	_, email, passhash := generateUserThenActivate()
	payload := service.AuthQuery{
		Email:    email,
		Passhash: passhash,
	}

	rv := reqPOST(endpoint.AuthorizeUser, payload)
	rMap := rv.Result.(map[string]interface{})
	return rMap["token"].(string)
}

func upgradeUser(token string) *models.Store {
	payload := service.BecomeAuctioneerQuery{
		Name:        faker.Company().Name(),
		Info:        faker.Company().CatchPhrase(),
		Province:    faker.Address().State(),
		Regency:     faker.Address().City(),
		SUBDistrict: faker.Address().StreetName(),
		Village:     faker.Address().StreetSuffix(),
		Address:     faker.Address().StreetAddress(),
	}
	rv := reqPOST(endpoint.BecomeAuctioneer, payload, token)
	rMap := rv.Result.(map[string]interface{})
	store := mapToJSON(rMap, &models.Store{}).(*models.Store)

	return store
}

func closeProduct(productID int64) {
	repo := repository.NewProductRepository()
	repo.CloseProduct(productID)
}

func createProduct(token string, storeID int64) (types.Product, error) {
	labels := []repository.LabelQuery{}
	labels = append(labels, repository.LabelQuery{
		Name:  "label_name",
		Value: "value",
	})
	payload := repository.NewProductQuery{
		StoreID:       storeID,
		ProductName:   faker.Commerce().ProductName(),
		ProductImages: []string{faker.Internet().Url()},
		Desc:          faker.RandomString(100),
		Condition:     1,
		ConditionAvg:  100,
		StartPrice:    50000,
		BidMultpl:     50000,
		ClosedAT:      utils.NOW.Add(time.Hour * 24).Format(time.RFC3339),
		Labels:        labels,
	}

	rv := reqPOST(endpoint.AddProduct, payload, token)
	if rv.Code != 0 {
		return types.Product{}, errors.New("Failed creating product")
	}
	product := types.Product{}
	resMap := rv.Result.(map[string]interface{})
	mapToJSON(resMap, &product)
	return product, nil
}

// convert map to json with param input is map interface,
// and output param is referenced interface
func mapToJSON(input map[string]interface{}, output interface{}) interface{} {
	config := &mapstructure.DecoderConfig{
		TagName: "json",
	}
	config.Result = output
	decoder, _ := mapstructure.NewDecoder(config)
	decoder.Decode(input)

	return output
}
