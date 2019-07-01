package handler

import (
	testSetup "Chefling_Go/tests"
	"testing"
	"net/http"
	"reflect"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"Chefling_Go/constant"
	"strings"
	"Chefling_Go/model"
	"Chefling_Go/authentication"
)

func TestSignUpHandler(t *testing.T) {
	validRequestBody := `{"firstName": "TEST","lastName": "XYZ","email": "test@gmail.com","password": "pass123!"}`
	invalidRequestBody := `{"firstName": TEST,"lastName": "XYZ","email": "test@gmail.com","password": "pass123!"}`
	validRequestBodyWithEmptyCredentials := `{"firstName": "TEST","lastName": "XYZ","email": "","password": ""}`

	userProfile := model.Profile{
		FirstName: "TEST",
		LastName:  "XYZ",
		Email:     testSetup.EMAIL,
		Password:  "pass123!",
	}

	testSetup.MockMongoClient.On("SaveUser", userProfile).Return(nil, nil)

	url := "/user/signup"

	tests := []struct {
		method             string
		requestBody        string
		expectedStatusCode int
		isUserExists       bool
		expectedBody       string
	}{
		{
			method:             constant.POST,
			requestBody:        validRequestBody,
			expectedStatusCode: http.StatusOK,
			isUserExists:       false,
			expectedBody:       authentication.GetToken(),
		},
		{
			method:             constant.POST,
			requestBody:        validRequestBodyWithEmptyCredentials,
			expectedStatusCode: http.StatusBadRequest,
			isUserExists:       true,
			expectedBody:       "Email or/&& Password is empty",
		},
		{
			method:             constant.POST,
			requestBody:        invalidRequestBody,
			expectedStatusCode: http.StatusBadRequest,
			isUserExists:       true,
			expectedBody:       "Invalid request body",
		},
		{
			method:             constant.GET,
			requestBody:        validRequestBody,
			expectedStatusCode: http.StatusMethodNotAllowed,
			isUserExists:       true,
			expectedBody:       ``,
		},
	}

	for testNumber, test := range tests {
		testSetup.MockMongoClient.On("IsUserExist", testSetup.EMAIL).Return(test.isUserExists, nil)

		request, _ := http.NewRequest(test.method, url, strings.NewReader(test.requestBody))
		signUpHandler := SignUpHandler(testSetup.MockMongoClient)
		rr := executeSignUpRequest(signUpHandler, request, url)

		if status := rr.Code; status != test.expectedStatusCode {
			t.Errorf(constant.TEST_STATUSCODE_MESSAGE, testNumber+1, status, test.expectedStatusCode, url)
		}
		res_body := rr.Body.String()
		if !reflect.DeepEqual(res_body, test.expectedBody) {
			t.Errorf(constant.TEST_RESPONSE_BODY_MESSSAGE, testNumber+1, res_body, test.expectedBody)
		}
	}
}

func executeSignUpRequest(handlerFunc http.HandlerFunc, request *http.Request, endpoint string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(endpoint, handlerFunc).Methods(constant.POST)
	router.ServeHTTP(rr, request)
	return rr
}
