package factory

import "Chefling_Go/model"

type DataStoreClient interface {
	SaveUser(user model.Profile)(interface{}, error)
	GetUserData(userKey string) (model.Profile, error)
	IsUserExist(userKey string) (bool, error)
	UpdateUser(user model.Profile) (error)
}

