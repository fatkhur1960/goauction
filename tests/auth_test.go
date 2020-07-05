package test

import (
	"testing"

	"github.com/fatkhur1960/goauction/app/service"
	"github.com/fatkhur1960/goauction/tests/endpoint"
	"github.com/go-playground/assert/v2"
)

func TestAuthorizeUser(t *testing.T) {
	_, email, passhash := generateUserThenActivate()
	payload := service.AuthQuery{
		Email:    email,
		Passhash: passhash,
	}

	rv := reqPOST(endpoint.AuthorizeUser, payload)
	assert.Equal(t, rv.Code, 0)

	rMap := rv.Result.(map[string]interface{})
	assert.NotEqual(t, rMap["token"], nil)
}

func TestUnauthorizeUser(t *testing.T) {
	token := authorizeUser()

	rv := reqPOST(endpoint.UnauthorizeUser, nil, token)
	assert.Equal(t, rv.Code, 0)

	rv1 := reqGET(endpoint.MeInfo, token)
	assert.Equal(t, rv1.Code, 4010)
}
