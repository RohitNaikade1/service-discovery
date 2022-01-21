package models

import "github.com/golang-jwt/jwt"

type User struct {
	ID         string `bson:"_id" json:"_id,omitempty"`
	First_Name string `json:"firstname,omitempty"`
	Last_Name  string `json:"lastname,omitempty"`
	UserName   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Email      string `json:"email,omitempty"`
	Role       string `json:"role,omitempty"`
	Created_At string `json:"created_at,omitempty"`
	Updated_At string `json:"updated_at,omitempty"`
}

type LoginDetails struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
