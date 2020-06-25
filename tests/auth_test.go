package test

import (
	"net/http"
	"testing"

	"github.com/fatkhur1960/goauction/app/service"
	"github.com/go-playground/assert/v2"
)

func TestAuthorizeUser(t *testing.T) {
	email, passhash := generateUserThenActivate()
	payload := service.AuthQuery{
		Email:    email,
		Passhash: passhash,
	}

	rv := reqPOST(AuthorizeUserEndpoint, payload)
	assert.Equal(t, rv.Code, 0)

	rMap := rv.Result.(map[string]interface{})
	assert.NotEqual(t, rMap["token"], nil)

	cleanUsers()
}

func TestUnauthorizeUser(t *testing.T) {
	token := authorizeUser()

	rv := reqPOST(UnauthorizeUserEndpoint, nil, token)
	assert.Equal(t, rv.Code, 0)

	rv1 := reqGET(MeInfoEndpoint, token)
	assert.Equal(t, rv1.Code, http.StatusUnauthorized)
}
