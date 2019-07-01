package handler

import (
	"net/http"
	"Chefling_Go/model"
	"Chefling_Go/factory"
	"Chefling_Go/util"
	"Chefling_Go/constant"
	"encoding/json"
)

func UpdateHandler(dataStoreClient factory.DataStoreClient) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer util.PanicHandler(writer, request)
		writer.Header().Set(constant.HEADER_CONTENT_TYPE, constant.HEADER_APPLICATION_JSON)
		var userProfile model.Profile
		decoder := json.NewDecoder(request.Body)
		invalidRequest := decoder.Decode(&userProfile)
		if invalidRequest == nil {
			email := util.GetUserID(request)
			userData := model.Profile{
				FirstName: userProfile.FirstName,
				LastName:  userProfile.LastName,
				Email:     email,
				Password:  userProfile.Password,
			}
			err := dataStoreClient.UpdateUser(userData)
			if err == nil {
				util.HandleSuccess(writer, `{"action" : "User Data Updated into Mongo "}`)
			} else {
				util.HandleError(writer, http.StatusBadRequest, `{"action" : "User Data Failed to Update in Mongo "}`)
			}
		} else {
			util.HandleError(writer, http.StatusBadRequest, "Invalid request body")
		}
		defer request.Body.Close()
	}
}
