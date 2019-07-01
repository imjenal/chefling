package handler

import (
	"net/http"
	"Chefling_Go/constant"
	"Chefling_Go/factory"
	"Chefling_Go/util"
	"Chefling_Go/model"
	"encoding/json"
	authenticate "Chefling_Go/authentication"
)

func SignInHandler(dataStoreClient factory.DataStoreClient) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer util.PanicHandler(writer, request)
		writer.Header().Set(constant.HEADER_CONTENT_TYPE, constant.HEADER_APPLICATION_JSON)
		var userCredentials model.UserCredentials
		decoder := json.NewDecoder(request.Body)
		invalidRequest := decoder.Decode(&userCredentials)
		if invalidRequest == nil {
			credentials, _ := dataStoreClient.GetUserData(userCredentials.Email)
			if util.IsEquals(credentials.Email, userCredentials.Email) && util.ComparePasswords(credentials.Password, userCredentials.Password) {
				util.HandleSuccess(writer, authenticate.GetToken())
			} else {
				util.HandleError(writer, http.StatusBadRequest, "Email Password doesn't match")
			}
		} else {
			util.HandleError(writer, http.StatusBadRequest, "Invalid request body")
		}

		defer request.Body.Close()
	}
}

