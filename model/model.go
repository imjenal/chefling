package model

type User struct {
	Id      string  `bson:"_id"`
	Profile Profile `json:"profile"`
}

type Profile struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type AccessToken struct {
	Token string `json:"auth_token"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
