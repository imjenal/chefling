package authentication

import (
	"net/http"
	"fmt"
	"strings"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"Chefling_Go/constant"
	"time"
	"github.com/dgrijalva/jwt-go"
	"encoding/json"
	"Chefling_Go/model"
	"Chefling_Go/util"
	"Chefling_Go/config"
)



func GetToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["admin"] = true
	claims["name"] = "testing"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	secretKey := []byte("secret")
	tokenString, _ := token.SignedString(secretKey)
	tokenObject := &model.AccessToken{
		Token: tokenString,
	}
	payload, _ := json.Marshal(tokenObject)

	return string(payload)
}

func VerifyAndServe(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := verify(request, config.GetAppConfiguration().SECRET_KEY)

		if err != nil {
			util.HandleError(writer, http.StatusUnauthorized, constant.UNAUTHORIZED)
			return
		}
		handler.ServeHTTP(writer, request)
	})
}


func verify(request *http.Request, secret string) error {
	if err := verifyDateUsedHeader(request); err != nil {
		return err
	}

	authorizationHeader := request.Header.Get(constant.HEADER_AUTHORIZATION)
	if authorizationHeader == "" {
		return fmt.Errorf("authorization header is not present")
	}

	_, signature, err := parseAuthorizationHeader(authorizationHeader)
	if err != nil {
		return err
	}

	if verifySignature(signature, generateCanonicalString(request), secret) {
		return nil
	}

	return fmt.Errorf("signature mismatch")
}

func verifySignature(signature, canonicalString, secret string) bool {
	expectedSignature := computeSignature(canonicalString, secret)
	return expectedSignature == signature
}

func generateCanonicalString(request *http.Request) string {
	uri := request.URL.EscapedPath()
	if uri == "" {
		uri = "/"
	}

	if request.URL.RawQuery != "" {
		uri = uri + constant.QUERY_SPLIT + request.URL.RawQuery
	}

	canonicalString := strings.Join([]string{
		request.Header.Get(constant.HEADER_CONTENT_TYPE),
		uri,
		request.Header.Get(constant.HEADER_DATE_USED),
	}, constant.DELIMITER)

	return canonicalString
}

func computeSignature(canonicalString, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(canonicalString))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature
}

func verifyDateUsedHeader(request *http.Request) error {
	date := request.Header.Get(constant.HEADER_DATE_USED)
	if date == "" {
		return fmt.Errorf("date header is not present")
	}
	requestTimeStamp, err := time.Parse(time.RFC1123, date)
	if err != nil {
		return fmt.Errorf("request timestamp is not in expected RFC1123 format")
	}

	currentTimeStamp := time.Now().UTC()
	if currentTimeStamp.After(requestTimeStamp.Add(15 * time.Minute)) {
		return fmt.Errorf("request has expired")
	}
	return nil
}

func parseAuthorizationHeader(authorizationHeader string) (id, signature string, err error) {
	var tokens []string

	if !strings.HasPrefix(authorizationHeader, "Authorization ") {
		goto malformed
	}

	tokens = strings.Split(strings.Split(authorizationHeader, "Authorization ")[1], constant.KEY_SEPARATOR)
	if len(tokens) != 2 || tokens[0] == "" || tokens[1] == "" {
		goto malformed
	}

	return tokens[0], tokens[1], nil

malformed:
	return "", "", fmt.Errorf("malformed Authorization Header: %s", authorizationHeader)
}

