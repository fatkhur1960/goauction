package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/router"
	"github.com/gin-gonic/gin"
)

var ts = httptest.NewServer(getTestingRoutes())

func getTestingRoutes() *gin.Engine {
	models.ConnectDatabaseTest()
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

	return result
}

func reqPOST(path string, data interface{}) app.Result {
	query, _ := json.Marshal(data)
	resp, err := http.Post(fmt.Sprintf("%s/api%s", ts.URL, path), "application/json", bytes.NewBuffer(query))
	return parseResult(resp, err)
}

func reqGET(path string) app.Result {
	resp, err := http.Get(path)
	return parseResult(resp, err)
}
