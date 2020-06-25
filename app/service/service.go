package service

import (
	"net/http"

	"github.com/fatkhur1960/goauction/app"
	"github.com/gin-gonic/gin"
)

type (
	// RegisterToken result after user register
	RegisterToken struct {
		Token string `json:"token"`
	}

	// EntriesResult for list result
	EntriesResult struct {
		Entries interface{} `json:"entries"`
		Count   int         `json:"count"`
	}

	// QueryEntries request type struct
	QueryEntries struct {
		Limit  int    `form:"limit"`
		Offset int    `form:"offset" binding:"required"`
		Query  string `form:"query"`
	}

	// IDQuery request type struct
	IDQuery struct {
		ID int64 `uri:"id" binding:"required"`
	}
)

const (
	typeJSON  = 1
	typeQuery = 2
	typeURI   = 3
)

func (q *IDQuery) validate(c *gin.Context, t int) {
	var err error
	switch t {
	case typeJSON:
		err = c.ShouldBindJSON(&q)
	case typeQuery:
		err = c.ShouldBindQuery(&q)
	case typeURI:
		err = c.ShouldBindUri(&q)
	}

	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	return
}

func (q *QueryEntries) validate(c *gin.Context, t int) {
	var err error
	switch t {
	case typeJSON:
		err = c.ShouldBindJSON(&q)
	case typeQuery:
		err = c.ShouldBindQuery(&q)
	case typeURI:
		err = c.ShouldBindUri(&q)
	}

	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	return
}

// Body request validation
func validateRequest(c *gin.Context, query interface{}) interface{} {
	if err := c.ShouldBindJSON(query); err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return nil
	}

	return query
}

// APIResult for json output
var APIResult = app.NewAPIResult()

// NoRouteHandler handle 404
func NoRouteHandler(c *gin.Context) {
	APIResult.Error(c, http.StatusNotFound, "Endpoint not found")
}
