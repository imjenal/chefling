package handler

import (
	testSetup "Chefling_Go/tests"
	"testing"
	"Chefling_Go/model"
	"Chefling_Go/constant"
	"net/http"
	"reflect"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"strings"
)

func TestSignInHandler(t *testing.T) {
	validRequestBody := `{"email": "test@gmail.com","password": "pass123!"}`
	invalidRequestBody := `{"email": "test@gmail.com","password":pass123!}`
	validRequestBodyWithWrongPassword := `{"email": "wrongtest@gmail.com","password": "wrong!"}`

	testSetup.MockMongoClient.On("GetUserData", testSetup.EMAIL).Return(model.Profile{
		FirstName: "TEST",
		LastName:  "XYZ",
		Email:     testSetup.EMAIL,
		Password:  "$2a$04$4rXSDVNFxu5PVGimsWmwUuClECk3.Dzt33uZ/JrIcNv/tYKUOS3xC",
	}, nil)

	testSetup.MockMongoClient.On("GetUserData", "wrongtest@gmail.com").Return(model.Profile{
		FirstName: "TEST",
		LastName:  "XYZ",
		Email:     "wrongtest@gmail.com",
		Password:  "wrongPassword",
	}, nil)


	url := "/user/signin"

	tests := []struct {
		method             string
		requestBody        string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			method:             constant.GET,
			requestBody:        validRequestBodyWithWrongPassword,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Email Password doesn't match",
		},
		{
			method:             constant.GET,
			requestBody:        invalidRequestBody,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Invalid request body",
		},
		{
			method:             constant.POST,
			requestBody:        validRequestBody,
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       ``,
		},
	}

	for testNumber, test := range tests {
		request, _ := http.NewRequest(test.method, url, strings.NewReader(test.requestBody))
		signInHandler := SignInHandler(testSetup.MockMongoClient)
		rr := executeSignInRequest(signInHandler, request, url)

		if status := rr.Code; status != test.expectedStatusCode {
			t.Errorf(constant.TEST_STATUSCODE_MESSAGE, testNumber+1, status, test.expectedStatusCode, url)
		}
		res_body := rr.Body.String()
		if !reflect.DeepEqual(res_body, test.expectedBody) {
			t.Errorf(constant.TEST_RESPONSE_BODY_MESSSAGE, testNumber+1, res_body, test.expectedBody)
		}
	}
}

func executeSignInRequest(handlerFunc http.HandlerFunc, request *http.Request, endpoint string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(endpoint, handlerFunc).Methods(constant.GET)
	router.ServeHTTP(rr, request)
	return rr
}
