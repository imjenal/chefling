package util

import (
	"fmt"
	"strings"
	"errors"
	"net/http"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HandleError(writer http.ResponseWriter, code int, message string) {
	writer.WriteHeader(code)
	writer.Write([]byte(message))
}

func HandleSuccess(writer http.ResponseWriter, message string) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(message))
}

func IsEmpty(data string) bool {
	return len(data) == 0
}

func IsEquals(input, testInput string) bool {
	return input == testInput
}

func GetUserID(request *http.Request) string {
	for k, userId := range mux.Vars(request) {
		if strings.EqualFold(k, "userId") {
			return userId
		}
	}
	return ""
}


func PanicHandler(w http.ResponseWriter, r *http.Request) {
	if r := recover(); r != nil {
		fmt.Sprintf("Recovered in f %v", r)
		var err error
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			err = errors.New("Unknown panic")
		}
		if err != nil {
			fmt.Printf("FAILED")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}


func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return false
	}
	return true

}