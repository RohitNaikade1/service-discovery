package models

import "github.com/golang-jwt/jwt"

type User struct {
	ID         string `bson:"_id" json:"_id,omitempty"`
	FirstName  string `bson:"firstname" json:"firstname,omitempty"`
	LastName   string `bson:"lastname" json:"lastname,omitempty"`
	UserName   string `bson:"username" json:"username,omitempty"`
	Password   string `bson:"password" json:"password,omitempty"`
	Email      string `bson:"email" json:"email,omitempty"`
	Role       string `bson:"role" json:"role,omitempty"`
	Created_At string `json:"created_at,omitempty"`
	Updated_At string `json:"updated_at,omitempty"`
}

type LoginDetails struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	ID       string `bson:"id" json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
