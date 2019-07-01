package main

import (
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"Chefling_Go/constant"
	"Chefling_Go/factory"
	"Chefling_Go/handler"
	"Chefling_Go/config"
	"Chefling_Go/authentication"
)

var AppConfig *config.Configuration = nil

func init() {
	AppConfig = config.GetAppConfiguration()
	if AppConfig.DATASTORE == "" || AppConfig.MONGO_SERVER == "" || AppConfig.MONGO_DB == "" {
		fmt.Println("Please set the Database")
	}

}

func main(){

	AppConfig = config.GetAppConfiguration()
	if (AppConfig.DATASTORE == "" || AppConfig.MONGO_SERVER == "" || AppConfig.MONGO_DB == "") {
		fmt.Println("Please set the Database credentials")
	}

	router := mux.NewRouter().StrictSlash(true)
	dataStoreClient := factory.GetDataStoreClient()

	allowedHeaders := handlers.AllowedHeaders([]string{constant.HEADER_ORIGIN, constant.HEADER_ACCEPT, constant.HEADER_CONTENT_TYPE, constant.HEADER_AUTHORIZATION, constant.HEADER_DATE_USED, constant.HEADER_X_REQUESTED_WITH})
	allowedOrigins := handlers.AllowedOrigins([]string{constant.ALL})
	allowedMethods := handlers.AllowedMethods([]string{constant.GET, constant.HEAD, constant.POST, constant.PUT, constant.DELETE, constant.OPTIONS})

	signUpHandler := handler.SignUpHandler(dataStoreClient)
	signInHandler := handler.SignInHandler(dataStoreClient)
	profileHandler := handler.ProfileHandler(dataStoreClient)
	updateHandler := handler.UpdateHandler(dataStoreClient)


	fileServer := http.FileServer(http.Dir("./swaggerui/"))
	router.PathPrefix("/api/").Handler(http.StripPrefix("/api/", fileServer))

	router.Handle("/user/signup", signUpHandler).Methods(constant.POST)
	router.Handle("/user/profile/{userId}", authentication.VerifyAndServe(profileHandler)).Methods(constant.GET)
	router.Handle("/user/profile/update/{userId}", authentication.VerifyAndServe(updateHandler)).Methods(constant.POST)
	router.Handle("/user/signin", signInHandler).Methods(constant.GET)

	fmt.Println("Application loaded successfully ")
	log.Fatal(http.ListenAndServe(":"+AppConfig.PORT, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router)))
}
