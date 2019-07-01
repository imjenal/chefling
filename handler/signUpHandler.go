package handler

import (
	"net/http"
	"encoding/json"
	"Chefling_Go/model"
	"Chefling_Go/constant"
	"Chefling_Go/factory"
	"Chefling_Go/util"
	authenticate "Chefling_Go/authentication"
)

func SignUpHandler(dataStoreClient factory.DataStoreClient) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer util.PanicHandler(writer, request)
		writer.Header().Set(constant.HEADER_CONTENT_TYPE, constant.HEADER_APPLICATION_JSON)

		var userData model.Profile
		decoder := json.NewDecoder(request.Body)
		invalidRequest := decoder.Decode(&userData)
		if invalidRequest == nil {
			if !util.IsEmpty(userData.Email) && !util.IsEmpty(userData.Password) {
				isUserExist,_ := dataStoreClient.IsUserExist(userData.Email)
				if !isUserExist  {
					dataStoreClient.SaveUser(userData)
					util.HandleSuccess(writer, authenticate.GetToken())
				} else {
					util.HandleError(writer, http.StatusBadRequest, "Email ID already exists")
				}
			} else {
				util.HandleError(writer, http.StatusBadRequest, "Email or/&& Password is empty")
			}
		} else {
			util.HandleError(writer, http.StatusBadRequest, "Invalid request body")
		}
		defer request.Body.Close()
	}
}
