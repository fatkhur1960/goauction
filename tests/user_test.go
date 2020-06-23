package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/service"
	"github.com/go-playground/assert/v2"
	"syreclabs.com/go/faker"
)

func TestRegisterUser(t *testing.T) {
	user := service.RegisterUserQuery{
		FullName: faker.Name().Name(),
		Email:    faker.Internet().Email(),
		PhoneNum: faker.PhoneNumber().CellPhone(),
	}
	rv := reqPOST("/user/v1/register", user)

	assert.Equal(t, rv.Code, 0)
}

func TestActivateUser(t *testing.T) {
	var userModel models.User
	userQuery := service.RegisterUserQuery{
		FullName: faker.Name().Name(),
		Email:    faker.Internet().Email(),
		PhoneNum: faker.PhoneNumber().CellPhone(),
	}

	rv := reqPOST("/user/v1/register", userQuery)
	assert.Equal(t, rv.Code, 0)
	assert.NotEqual(t, rv.Result, nil)

	activate := service.ActivateUserQuery{
		Token:    fmt.Sprintf("%v", rv.Result),
		Passhash: "123123",
	}
	rv1 := reqPOST("/user/v1/activate", activate)
	assert.Equal(t, rv1.Code, 0)

	userData, _ := json.Marshal(rv1.Result)
	json.Unmarshal(userData, &userModel)
	assert.Equal(t, userQuery.Email, userModel.Email)
}
