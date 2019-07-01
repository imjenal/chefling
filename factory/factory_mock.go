package factory

import (
	"github.com/stretchr/testify/mock"
	"Chefling_Go/model"
)

func SetDatastoreClient(client DataStoreClient) {
	dataStoreClient = client
}

type MongoMockClientImpl struct {
	mock.Mock
	mongoServer string
}

func (client *MongoMockClientImpl) SaveUser(user model.Profile) (interface{}, error) {
	args := client.Called(user)
	return args.Get(0), args.Error(1)
}

func (client *MongoMockClientImpl) IsUserExist(userKey string) (bool, error) {
	args := client.Called(userKey)
	return args.Get(0).(bool), args.Error(1)
}

func (client *MongoMockClientImpl) GetUserData(userKey string) (model.Profile, error) {
	args := client.Called(userKey)
	return args.Get(0).(model.Profile), args.Error(1)
}

func (client *MongoMockClientImpl) UpdateUser(user model.Profile) (error)  {
	args := client.Called(user)
	return args.Error(1)
}
