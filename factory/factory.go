package factory

import (
	"Chefling_Go/constant"
	"gopkg.in/mgo.v2"
	"Chefling_Go/config"
)

var dataStoreClient DataStoreClient = nil

func GetDataStoreClient() DataStoreClient {
	if dataStoreClient == nil {
		dataStore := config.GetAppConfiguration().DATASTORE
		if dataStore == constant.MONGO {
			SetMongoClient(config.GetAppConfiguration().MONGO_SERVER, nil)
		}
	}
	return dataStoreClient
}


func SetMongoClient(server string, previousSession *mgo.Session) {
	if previousSession != nil {
		defer previousSession.Close()
	}
	session, err := mgo.Dial(server)
	if err == nil {
		dataStoreClient = &MongoClientImpl{mongoServer: server, session: session}
	}
}
