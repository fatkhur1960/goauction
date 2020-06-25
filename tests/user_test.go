package test

import (
	"fmt"
	"testing"

	"github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/app/service"
	"github.com/go-playground/assert/v2"
	"syreclabs.com/go/faker"
)

func TestRegisterUser(t *testing.T) {
	u := service.RegisterUserQuery{
		FullName: faker.Name().Name(),
		Email:    faker.Internet().Email(),
		PhoneNum: faker.PhoneNumber().CellPhone(),
	}
	rv := reqPOST(RegisterUserEndpoint, u)
	assert.Equal(t, rv.Code, 0)

	rMap := rv.Result.(map[string]interface{})
	assert.Equal(t, rMap["email"], u.Email)
}

func TestActivateUser(t *testing.T) {
	u := service.RegisterUserQuery{
		FullName: faker.Name().Name(),
		Email:    faker.Internet().Email(),
		PhoneNum: faker.PhoneNumber().CellPhone(),
	}

	rv := reqPOST(RegisterUserEndpoint, u)
	assert.Equal(t, rv.Code, 0)

	rMap := rv.Result.(map[string]interface{})
	assert.Equal(t, rMap["email"], u.Email)

	activate := service.ActivateUserQuery{
		Token:    fmt.Sprintf("%v", rMap["token"]),
		Passhash: faker.Internet().Password(8, 12),
	}
	rv1 := reqPOST(ActivateUserEndpoint, activate)
	assert.Equal(t, rv1.Code, 0)

	rMap2 := rv.Result.(map[string]interface{})
	assert.Equal(t, u.Email, rMap2["email"])

	cleanUsers()
}

func TestUpdateUser(t *testing.T) {
	token := authorizeUser()
	u := repository.UpdateUserQuery{
		FullName: faker.Name().Name(),
		Email:    faker.Internet().Email(),
		PhoneNum: faker.PhoneNumber().CellPhone(),
		Address:  faker.Address().String(),
		Avatar:   faker.Internet().Url() + ".jpg",
	}

	rv := reqPOST(UpdateUserInfoEndpoint, u, token)
	assert.Equal(t, rv.Code, 0)

	rMap := rv.Result.(map[string]interface{})
	assert.Equal(t, rMap["full_name"], u.FullName)
	assert.Equal(t, rMap["email"], u.Email)
	assert.Equal(t, rMap["phone_num"], u.PhoneNum)
	assert.Equal(t, rMap["address"], u.Address)
	assert.Equal(t, rMap["avatar"], u.Avatar)

	cleanUsers()
}
