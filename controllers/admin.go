package controllers

import (
	"net/http"
	"service-discovery/env"
	log "service-discovery/middlewares"
	"service-discovery/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var Logger = log.Logger()

func GenerateAdminToken(c *gin.Context) {
	expirationTime := time.Now().Add(time.Minute * 30)

	claims := &models.Claims{
		Username: env.GetEnvironmentVariable("ADMIN_USERNAME"),
		Password: env.GetEnvironmentVariable("ADMIN_PASSWORD"),
		Role:     "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		Logger.Error(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
