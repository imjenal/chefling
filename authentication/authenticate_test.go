package authentication

import (
	testSetup "Chefling_Go/tests"
	"testing"
	"net/http"
	"time"
	"errors"
	"github.com/stretchr/testify/assert"
	"Chefling_Go/constant"
	"reflect"
	"encoding/json"
	"Chefling_Go/model"
)

func TestGetToken(t *testing.T) {
	tokenObject := model.AccessToken{Token: testSetup.TOKEN}
	expectedInBytes, _ := json.Marshal(&tokenObject)
	expectedToken := string(expectedInBytes)
	actualToken := GetToken()
	if reflect.TypeOf(actualToken) != reflect.TypeOf(expectedToken) {
		t.Errorf("Token does not match: Expected: %v but Got: %v", expectedToken, actualToken)
	}
}

func TestVerify(t *testing.T) {
	request, _ := http.NewRequest(constant.POST, "http://example.coma/some/path?x=1&b=2", nil)
	request.Header.Set(constant.HEADER_DATE_USED, "Sun, 30 Jun 2020 19:34:03 GMT")

	tests := []struct {
		authorizationHeader string
		expected            error
	}{
		{
			authorizationHeader: "Authorization me:WrongSignature",
			expected:            errors.New("signature mismatch"),
		},
	}

	for _, test := range tests {
		request.Header.Set(constant.HEADER_AUTHORIZATION, test.authorizationHeader)
		actual := verify(request, testSetup.SECRET_KEY)
		assert.Equal(t, test.expected, actual)
	}

}

func TestVerify_NoAuthorizationHeader(t *testing.T) {
	request, _ := http.NewRequest(constant.POST, "http://example.com", nil)
	date := time.Now().Add(10 * time.Minute).Format(time.RFC1123)
	request.Header.Set(constant.HEADER_DATE_USED, date)
	expected := errors.New("authorization header is not present")
	actual := verify(request, "secret")
	assert.Equal(t, expected, actual)
}

func TestParseAuthorizationHeader(t *testing.T) {
	_, _, actualError := parseAuthorizationHeader("NotAuthorization here")
	expectedError := errors.New("malformed Authorization Header: NotAuthorization here")
	assert.Equal(t, expectedError, actualError)

	_, _, actualError = parseAuthorizationHeader("Authorization noseparator")
	expectedError = errors.New("malformed Authorization Header: Authorization noseparator")
	assert.Equal(t, expectedError, actualError)

	_, _, actualError = parseAuthorizationHeader("Authorization nosig:")
	expectedError = errors.New("malformed Authorization Header: Authorization nosig:")
	assert.Equal(t, expectedError, actualError)

	_, _, actualError = parseAuthorizationHeader("Authorization :noaccessID")
	expectedError = errors.New("malformed Authorization Header: Authorization :noaccessID")
	assert.Equal(t, expectedError, actualError)

	id, sig, actualError := parseAuthorizationHeader("Authorization me:sig")
	assert.Equal(t, "me", id)
	assert.Equal(t, "sig", sig)
}

func TestVerifyDateUsedHeader_NoDateUsedHeader(t *testing.T) {
	request, _ := http.NewRequest(constant.POST, "http://example.com", nil)
	expected := errors.New("date header is not present")
	actual := verifyDateUsedHeader(request)
	assert.Equal(t, expected, actual)
}

func TestVerifyDateUsedHeader_DateNotInFormat(t *testing.T) {
	request, _ := http.NewRequest(constant.POST, "http://example.com", nil)
	request.Header.Set(constant.HEADER_DATE_USED, "Sunday, 30 Jun 2019 05:20:40 IST")
	expected := errors.New("request timestamp is not in expected RFC1123 format")
	actual := verifyDateUsedHeader(request)
	assert.Equal(t, expected, actual)
}

func TestVerifyDateUsedHeader_RequestExpired(t *testing.T) {
	request, _ := http.NewRequest(constant.POST, "http://example.com", nil)
	date := time.Now().Add(- 20 * time.Minute).Format(time.RFC1123)
	request.Header.Set(constant.HEADER_DATE_USED, date)
	expected := errors.New("request has expired")
	actual := verifyDateUsedHeader(request)
	assert.Equal(t, expected, actual)
}

func TestVerifySignature(t *testing.T) {
	request, _ := http.NewRequest(constant.POST, "http://example.com", nil)
	request.Header.Set(constant.HEADER_DATE_USED, "Sun, 30 Jun 2019 19:34:03 GMT")
	request.Header.Set(constant.HEADER_AUTHORIZATION, "Authorization me:UStchMaZA4pF5nrqadDP/JbBoCg=")
	canonicalString := "application/json,WnNni3tnQAUFZDSkgFRwfQ==,/a?b=c,Sun, 30 Jun 2019 19:34:03 GMT"

	tests := []struct {
		signature string
		expected  bool
	}{
		{
			signature: "0QbvvUQ7O9GxSOIyudh1r52ky5w=",
			expected:  true,
		},
		{
			signature: "WrongSignature",
			expected:  false,
		}}

	for testNumber, test := range tests {
		actual := verifySignature(test.signature, canonicalString, "secret")
		if actual != test.expected {
			t.Errorf("Test No %d: Expected : %v but Got: %v", testNumber+1, test.expected, actual)
		}
	}

}

func TestGenerateCanonicalString(t *testing.T) {
	date := time.Now().Add(10 * time.Minute).Format(time.RFC1123)
	request, _ := http.NewRequest(constant.POST, "http://example.com", nil)
	assert.Equal(t, ",/,", generateCanonicalString(request))

	request.URL.Path = "/"
	assert.Equal(t, ",/,", generateCanonicalString(request))

	request.URL.Path = "/some/path"
	assert.Equal(t, ",/some/path,", generateCanonicalString(request))

	request.URL.RawQuery = "x=1&b=2"
	assert.Equal(t, ",/some/path?x=1&b=2,", generateCanonicalString(request))

	request.Header.Add(constant.HEADER_CONTENT_TYPE, "application/json")
	request.Header.Add(constant.HEADER_DATE_USED, date)

	expectedCanonicalString := "application/json,/some/path?x=1&b=2," + date
	assert.Equal(t, expectedCanonicalString, generateCanonicalString(request))
}

func TestComputeSignature(t *testing.T) {
	canonicalString := "application/json,WnNni3tnQAUFZDSkgFRwfQ==,/a?b=c,Sun, 30 Jun 2019 19:34:03 GMT"
	actual := computeSignature(canonicalString, "secret")
	expected := "0QbvvUQ7O9GxSOIyudh1r52ky5w="
	if actual != expected {
		t.Errorf("Expected: %v but Got: %v", expected, actual)
	}
}
