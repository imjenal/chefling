package handler

import (
	testSetup "Chefling_Go/tests"
	"testing"
	"Chefling_Go/model"
	"Chefling_Go/constant"
	"net/http"
	"strings"
	"reflect"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"fmt"
)

func TestUpdateHandler(t *testing.T) {
	validRequestBody := `{"firstName": "TEST","lastName": "XYZ","email": "test@gmail.com","password": "pass123!"}`
	invalidRequestBody := `{"firstName": TEST,"lastName": "XYZ","email": "test@gmail.com","password": "pass123!"}`

	userProfile := model.Profile{
		FirstName: "TEST",
		LastName:  "XYZ",
		Email:     testSetup.EMAIL,
		Password:  "pass123!",
	}

	testSetup.MockMongoClient.On("UpdateUser", userProfile).Return(nil, nil)

	tests := []struct {
		url                string
		method             string
		requestBody        string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			url:                fmt.Sprintf("/user/profile/update/%s", testSetup.EMAIL),
			method:             constant.POST,
			requestBody:        validRequestBody,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"action" : "User Data Updated into Mongo "}`,
		},
		{
			url:                fmt.Sprintf("/user/profile/update/%s", testSetup.EMAIL),
			method:             constant.POST,
			requestBody:        invalidRequestBody,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Invalid request body",
		},
		{
			url:                fmt.Sprintf("/user/profile/update/%s", testSetup.EMAIL),
			method:             constant.GET,
			requestBody:        validRequestBody,
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       ``,
		},
	}

	for testNumber, test := range tests {
		request, _ := http.NewRequest(test.method, test.url, strings.NewReader(test.requestBody))
		updateHandler := UpdateHandler(testSetup.MockMongoClient)
		rr := executeUpdateRequest(updateHandler, request, "/user/profile/update/{userId}")

		if status := rr.Code; status != test.expectedStatusCode {
			t.Errorf(constant.TEST_STATUSCODE_MESSAGE, testNumber+1, status, test.expectedStatusCode, test.url)
		}
		res_body := rr.Body.String()
		if !reflect.DeepEqual(res_body, test.expectedBody) {
			t.Errorf(constant.TEST_RESPONSE_BODY_MESSSAGE, testNumber+1, res_body, test.expectedBody)
		}
	}
}

func executeUpdateRequest(handlerFunc http.HandlerFunc, request *http.Request, endpoint string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(endpoint, handlerFunc).Methods(constant.POST)
	router.ServeHTTP(rr, request)
	return rr
}
