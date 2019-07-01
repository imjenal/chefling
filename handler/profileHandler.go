package handler

import (
	"net/http"
	"encoding/json"
	"Chefling_Go/factory"
	"Chefling_Go/util"
	"Chefling_Go/constant"
)

func ProfileHandler(dataStoreClient factory.DataStoreClient) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer util.PanicHandler(writer, request)
		writer.Header().Set(constant.HEADER_CONTENT_TYPE, constant.HEADER_APPLICATION_JSON)
		email := util.GetUserID(request)
		userData, _ := dataStoreClient.GetUserData(email)
		jsonResult, err := json.Marshal(userData)
		if err == nil {
			util.HandleSuccess(writer, string(jsonResult))
		} else {
			util.HandleError(writer, http.StatusBadRequest, "Userdata Json Parse Error")
		}
	}
}
