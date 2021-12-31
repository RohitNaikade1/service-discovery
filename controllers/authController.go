package controllers

import (
	"net/http"
	log "service-discovery/middlewares"
	"service-discovery/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	jwtKey = []byte("secret_key")
)

func Authenticate(c *gin.Context) {
	bearer := c.Request.Header.Get("Authorization")
	if bearer == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization header provided"})
		c.Abort()
	} else {
		split := strings.Split(bearer, "Bearer ")
		if len(split) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Bearer token. Please add Bearer + token in Authorization"})
			c.Abort()
			return
		}
		token := split[1]

		claims, err := ValidateToken(token)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
		} else {
			c.Set("username", claims.Username)
			c.Set("password", claims.Password)
			c.Set("role", claims.Role)
			c.Next()
		}
	}

}

func ValidateToken(signedToken string) (claim *models.Claims, result string) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		log.Logger().Error(err.Error())
		result = err.Error()
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		result = "token is invalid"
	}

	if claims.ExpiresAt < time.Now().Unix() {
		result = "token has expired"
	}
	return claims, result
}
