package factory

import (
	"fmt"
	"strings"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"Chefling_Go/model"
	"Chefling_Go/constant"
	"Chefling_Go/util"
	"Chefling_Go/config"
)

type MongoClientImpl struct {
	mongoServer string
	session     *mgo.Session
}

func (client *MongoClientImpl) SaveUser(user model.Profile) (interface{}, error) {
	collection := fmt.Sprintf(constant.COLLECTION_USERS)
	userKey := user.Email
	session := client.session.Copy()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(config.GetAppConfiguration().MONGO_DB).C(collection)

	user.Password = util.HashAndSalt([]byte(user.Password))
	var userData = model.User{
		Id:      user.Email,
		Profile: user,
	}

	_, err := c.UpsertId(userKey, userData)

	if err != nil {
		msg := fmt.Sprintf("[MONGO] Error Saving User: %s in collection %s, error : %v", user.Email, collection, err)
		fmt.Println(msg)
		return nil, err
	}
	return nil, nil
}

func (client *MongoClientImpl) IsUserExist(userKey string) (bool, error) {
	collection := fmt.Sprintf(constant.COLLECTION_USERS)
	session := client.session.Copy()
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(config.GetAppConfiguration().MONGO_DB).C(collection)
	var result interface{}
	err := c.Find(bson.M{constant.MONGO_ID: userKey}).One(&result)
	if err != nil {
		if strings.Contains(err.Error(), constant.MONGO_ERROR_NOT_FOUND) {
			return false, err
		}
		msg := fmt.Sprintf("[MONGO] Error retrieving the User: %s in collection: %s, error : %v", userKey, collection, err)
		fmt.Println(msg)
	}
	return true, nil

}

func (client *MongoClientImpl) GetUserData(userKey string) (model.Profile, error) {
	collection := fmt.Sprintf(constant.COLLECTION_USERS)
	session := client.session.Copy()
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(config.GetAppConfiguration().MONGO_DB).C(collection)
	result := model.User{}
	err := c.Find(bson.M{constant.MONGO_ID: userKey}).Select(bson.M{"profile": 1}).One(&result)

	if err != nil {
		if strings.Contains(err.Error(), constant.MONGO_ERROR_NOT_FOUND) {
			msg := fmt.Sprintf("[MONGO] User: %s doesn't exist in collection: %s", userKey, collection)
			fmt.Println(msg)
		} else {
			msg := fmt.Sprintf("[MONGO] Error retrieving the User: %s in collection: %s, error : %v", userKey, collection, err)
			fmt.Println(msg)
		}
		return model.Profile{}, err
	}

	return result.Profile, nil
}

func (client *MongoClientImpl) UpdateUser(user model.Profile) ( error) {
	collection := fmt.Sprintf(constant.COLLECTION_USERS)
	session := client.session.Copy()
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(config.GetAppConfiguration().MONGO_DB).C(collection)

	user.Password = util.HashAndSalt([]byte(user.Password))
	data := model.Profile{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}

	selector := bson.M{constant.MONGO_ID: user.Email}
	changes := bson.M{"$set": bson.M{"profile": data}}

	err := c.Update(selector, changes)

	if err != nil {
		msg := fmt.Sprintf("[MONGO] Error Updating User: %s in collection: for %s, error : %v", user.Email, collection, err)
		fmt.Println(msg)
		return  err
	}
	return nil
}
