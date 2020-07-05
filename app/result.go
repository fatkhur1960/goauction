package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	// APIResult output
	APIResult interface {
		Success(c *gin.Context, res interface{})
		Error(c *gin.Context, code int, description string)
	}

	// Result untuk json response
	Result struct {
		Code        int         `json:"code"`
		Description string      `json:"description,omitempty"`
		Result      interface{} `json:"result"`
	}
)

// NewAPIResult instance for APIResult
func NewAPIResult() APIResult {
	return &Result{
		Code:        0,
		Description: "",
		Result:      nil,
	}
}

// Success api result
func (r *Result) Success(c *gin.Context, res interface{}) {
	var output map[string]interface{}

	r.Code = 0
	r.Result = res
	r.Description = ""
	data, _ := json.Marshal(r)

	json.Unmarshal(data, &output)
	c.JSON(http.StatusOK, output)
}

// Error api result
func (r *Result) Error(c *gin.Context, code int, description string) {
	var output map[string]interface{}

	code, _ = strconv.Atoi(fmt.Sprintf("%d0", code))
	r.Code = code
	r.Description = description
	r.Result = nil
	data, _ := json.Marshal(r)

	json.Unmarshal(data, &output)
	c.JSON(0, output)
	c.Abort()
}
