package test

import (
	"testing"

	"github.com/fatkhur1960/goauction/app/service"
	"github.com/fatkhur1960/goauction/tests/endpoint"
	"github.com/stretchr/testify/assert"
)

func TestCreateChatRoom(t *testing.T) {
	token := authorizeUser()
	userID, _, _ := generateUserThenActivate()
	payload := service.CreateChatQuery{
		UserID: userID,
	}

	rv := reqPOST(endpoint.CreateChatRoom, payload, token)
	assert.Equal(t, 0, rv.Code)
}

func TestCreateChatRoomUnauthorized(t *testing.T) {
	token := ""
	userID, _, _ := generateUserThenActivate()
	payload := service.CreateChatQuery{
		UserID: userID,
	}

	rv := reqPOST(endpoint.CreateChatRoom, payload, token)
	assert.Equal(t, 4010, rv.Code)
}

func TestListChatRooms(t *testing.T) {
	token := authorizeUser()
	userID, _, _ := generateUserThenActivate()
	payload := service.CreateChatQuery{
		UserID: userID,
	}

	reqPOST(endpoint.CreateChatRoom, payload, token)

	rv := reqGET(endpoint.ListChatRooms+"?offset=0&limit=10", token)
	entries := service.EntriesResult{}
	mapToJSON(rv.Result.(map[string]interface{}), &entries)
	assert.NotEqual(t, 0, entries.Count)
}
