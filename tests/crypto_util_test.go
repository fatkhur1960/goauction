package test

import (
	"log"
	"testing"

	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/go-playground/assert/v2"
	"syreclabs.com/go/faker"
)

func TestPasswordHashing(t *testing.T) {
	passhash := faker.Internet().Password(6, 8)
	hash, err := utils.GeneratePasshash(passhash)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, utils.CheckPasshash(passhash, hash), true)
}
