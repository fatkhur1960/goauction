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
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"syreclabs.com/go/faker"
)

var (
	ts     = httptest.NewServer(getTestingRoutes())
	client = &http.Client{}
)

func getTestingRoutes() *gin.Engine {
	models.ConnectDatabaseTest()
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

	query, _ := json.MarshalIndent(payload, "", "\t")
	if os.Getenv("TEST_LOG") == "debug" {
		log.Printf("=> sending payload to %s %s\n", url, string(query))
	}
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(query),
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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

func generateUserThenActivate() (string, string) {
	u := service.RegisterUserQuery{
		FullName: faker.Name().Name(),
		Email:    faker.Internet().Email(),
		PhoneNum: faker.PhoneNumber().CellPhone(),
	}

	passhash := faker.Internet().Password(8, 12)

	rv := reqPOST(RegisterUserEndpoint, u)
	rMap := rv.Result.(map[string]interface{})

	activate := service.ActivateUserQuery{
		Token:    fmt.Sprintf("%v", rMap["token"]),
		Passhash: passhash,
	}
	reqPOST(ActivateUserEndpoint, activate)

	return u.Email, passhash
}

func authorizeUser() string {
	email, passhash := generateUserThenActivate()
	payload := service.AuthQuery{
		Email:    email,
		Passhash: passhash,
	}

	rv := reqPOST(AuthorizeUserEndpoint, payload)
	rMap := rv.Result.(map[string]interface{})
	return rMap["token"].(string)
}

func createProduct(token string) (interface{}, error) {
	payload := repository.NewProductQuery{
		ProductName:   faker.Commerce().ProductName(),
		ProductImages: []string{faker.Internet().Url()},
		Desc:          faker.RandomString(100),
		Condition:     1,
		ConditionAvg:  100,
		StartPrice:    float64(faker.Commerce().Price()),
		BidMultpl:     float64(faker.Commerce().Price()),
		ClosedAT:      utils.NOW.Add(time.Hour * 24).Format(time.RFC3339),
		Labels:        []string{faker.RandomString(10)},
	}

	rv := reqPOST(AddProductEndpoint, payload, token)
	if rv.Code != 0 {
		return nil, errors.New("Failed creating product")
	}
	product := service.Product{}
	resMap := rv.Result.(map[string]interface{})
	mapstructure.Decode(resMap, &product)

	return product, nil
}

func cleanUsers() {
	userRepo := repository.UserRepository{
		UserQs: models.NewUserQuerySet(models.DB),
	}

	productRepo := repository.ProductRepository{
		ProductQs: models.NewProductQuerySet(models.DB),
	}

	userRepo.CleanUpUser()
	productRepo.CleanUpProduct()
}
