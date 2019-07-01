package handler

import (
	testSetup "Chefling_Go/tests"
	"testing"
	"net/http"
	"reflect"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"Chefling_Go/constant"
	"fmt"
	"Chefling_Go/model"
)

func TestProfileHandler(t *testing.T) {
	testSetup.MockMongoClient.On("GetUserData", testSetup.EMAIL).Return(model.Profile{
		FirstName: "TEST",
		LastName: "XYZ",
		Email: testSetup.EMAIL,
		Password: "pass123!",
	}, nil)

	tests := []struct {
		url                string
		method             string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			url:                fmt.Sprintf("/user/profile/%s", testSetup.EMAIL),
			method:             constant.GET,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"firstName":"TEST","lastName":"XYZ","email":"test@gmail.com","password":"pass123!"}`,
		},
		{
			url:                fmt.Sprintf("/user/profile/%s", testSetup.EMAIL),
			method:             constant.POST,
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       ``,
		},
	}

	for testNumber, test := range tests {
		request, _ := http.NewRequest(test.method, test.url,  nil)
		profileHandler := ProfileHandler(testSetup.MockMongoClient)
		rr := executeProfileRequest(profileHandler, request, "/user/profile/{userId}")

		if status := rr.Code; status != test.expectedStatusCode {
			t.Errorf(constant.TEST_STATUSCODE_MESSAGE, testNumber+1, status, test.expectedStatusCode, test.url)
		}
		res_body := rr.Body.String()
		if !reflect.DeepEqual(res_body, test.expectedBody) {
			t.Errorf(constant.TEST_RESPONSE_BODY_MESSSAGE, testNumber+1, res_body, test.expectedBody)
		}
	}
}

func executeProfileRequest(handlerFunc http.HandlerFunc, request *http.Request, endpoint string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(endpoint, handlerFunc).Methods(constant.GET)
	router.ServeHTTP(rr, request)
	return rr
}
