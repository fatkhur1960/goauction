package service

import (
	"net/http"

	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/types"
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
		Limit  int    `form:"limit" binding:"required"`
		Offset int    `form:"offset"`
		Query  string `form:"query"`
		Filter string `form:"filter"`
	}

	// QueryMessages request type struct
	QueryMessages struct {
		ChatID int64  `form:"chat_id" binding:"required"`
		Limit  int    `form:"limit" binding:"required"`
		Offset int    `form:"offset"`
		Query  string `form:"query"`
		Filter string `form:"filter"`
	}

	// QueryProducts request type struct
	QueryProducts struct {
		ProductID int64  `form:"product_id" binding:"required"`
		Limit     int    `form:"limit" binding:"required"`
		Offset    int    `form:"offset"`
		Query     string `form:"query"`
		Filter    string `form:"filter"`
	}

	// IDQuery request type struct
	IDQuery struct {
		ID int64 `form:"id" binding:"required"`
	}
)

func (q *IDQuery) validate(c *gin.Context, t types.ValidatorType) error {
	var err error
	switch t {
	case types.ValidateJSON:
		err = c.ShouldBindJSON(&q)
	case types.ValidateQuery:
		err = c.ShouldBindQuery(&q)
	case types.ValidateURI:
		err = c.ShouldBindUri(&q)
	}

	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return err
	}
	return nil
}

func (q *QueryEntries) validate(c *gin.Context, t types.ValidatorType) error {
	var err error
	switch t {
	case types.ValidateJSON:
		err = c.ShouldBindJSON(&q)
	case types.ValidateQuery:
		err = c.ShouldBindQuery(&q)
	case types.ValidateURI:
		err = c.ShouldBindUri(&q)
	}

	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

// Body request validation
func validateRequest(c *gin.Context, query interface{}) error {
	if err := c.ShouldBindJSON(query); err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

// APIResult for json output
var APIResult = app.NewAPIResult()

// NoRouteHandler handle 404
func NoRouteHandler(c *gin.Context) {
	APIResult.Error(c, http.StatusNotFound, "Endpoint not found")
}
