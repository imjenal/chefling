package tests

import (
	"Chefling_Go/factory"
)

var MockMongoClient = new(factory.MongoMockClientImpl)

const EMAIL = "test@gmail.com"
const TOKEN = "token"
const SECRET_KEY = "secret"

func init() {
	factory.SetDatastoreClient(MockMongoClient)
}

