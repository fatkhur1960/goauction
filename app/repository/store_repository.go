package repository

import (
	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/utils"
)

type (
	// StoreRepository init implementation
	StoreRepository struct {
		userQs  models.UserQuerySet
		storeQs models.StoreQuerySet
	}
)

// NewStoreRepository intance
func NewStoreRepository() *StoreRepository {
	return &StoreRepository{
		userQs:  models.NewUserQuerySet(app.DB),
		storeQs: models.NewStoreQuerySet(app.DB),
	}
}

// CreateStore dao for creating store after user upgraded
func (s *StoreRepository) CreateStore(ownerID int64, name string, info string, province string, regency string, subDistrict string, village string, address string) (models.Store, error) {
	store := models.Store{
		OwnerID:     ownerID,
		Name:        name,
		Info:        info,
		Province:    province,
		Regency:     regency,
		SUBDistrict: subDistrict,
		Village:     village,
		Address:     address,
		LastUpdated: &utils.NOW,
		TS:          &utils.NOW,
	}

	if err := store.Create(app.DB); err != nil {
		return models.Store{}, err
	}

	s.userQs.IDEq(ownerID).GetUpdater().SetType(2).Update()

	return store, nil
}

// GetByID get store by id
func (s *StoreRepository) GetByID(storeID int64) (models.Store, error) {
	store := models.Store{}
	if err := s.storeQs.IDEq(storeID).One(&store); err != nil {
		return store, err
	}

	return store, nil
}

// GetStoreByOwnerID get store by id
func (s *StoreRepository) GetStoreByOwnerID(ownerID int64) (models.Store, error) {
	store := models.Store{}
	if err := s.storeQs.OwnerIDEq(ownerID).One(&store); err != nil {
		return store, err
	}

	return store, nil
}
