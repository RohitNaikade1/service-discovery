package controllers

import (
	"net/http"
	"service-discovery/database"
	"service-discovery/env"
	"service-discovery/helpers"

	"service-discovery/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username := c.GetString("username")
	userPassword := c.GetString("password")
	role := c.GetString("role")

	sysAdmin := VerifyParentAdmin(username, userPassword, role)
	appUser := GetCurrentLoggedInUser(username, userPassword, role)

	if sysAdmin || appUser.Role == "admin" {
		var user models.User
		user.ID = primitive.NewObjectID().Hex()
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid json provided"})
		}
		password := user.Password
		hashPassword, err := helpers.HashPassword(password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		user.Password = hashPassword
		user.Created_At = time.Now().Local().String()
		user.Updated_At = time.Now().Local().String()

		if database.ValidateDocument(env.USER_COLLECTION, bson.M{"email": user.Email, "username": user.UserName}) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username or email already exists"})
		} else {
			result := database.Insert(env.USER_COLLECTION, user)
			c.JSON(http.StatusOK, gin.H{"status": "inserted", "id": result.InsertedID})
		}
	} else {
		c.JSON(http.StatusUnauthorized, "Not authorized to add users")
	}
	Logger.Debug("FUNCEXIT")
}

func Login(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	c.Header("Access-Control-Allow-Origin", "*")
	var login models.LoginDetails
	var user models.User
	err := c.BindJSON(&login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	//collection := database.UserCollection()
	if database.ValidateDocument(env.USER_COLLECTION, bson.M{"username": login.UserName}) {
		//collection.FindOne(context.Background(), bson.M{"username": login.UserName}).Decode(&user)
		database.Read(env.USER_COLLECTION, bson.M{"username": login.UserName}).Decode(&user)
		match := helpers.CheckPasswordHash(login.Password, user.Password)

		if match {

			expirationTime := time.Now().Add(time.Minute * 30)
			claims := &models.Claims{
				Username: user.UserName,
				Password: user.Password,
				Role:     user.Role,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			helpers.PrintError(err)
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect password"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
	}
	Logger.Debug("FUNCEXIT")
}

func GetUsers(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")
	var arr []string
	if database.ValidateCollection(env.USER_COLLECTION) {
		results := database.ReadAll(env.USER_COLLECTION)
		for _, users := range results {
			response := helpers.Encode(users)
			arr = append(arr, string(response))
		}
		stringByte := "[" + strings.Join(arr, " ,") + "]"
		sysAdmin := VerifyParentAdmin(username, password, role)
		appUser := GetCurrentLoggedInUser(username, password, role)
		if sysAdmin || appUser.Role == "admin" {
			c.Data(http.StatusOK, "application/json", []byte(stringByte))
		} else {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
	}
	Logger.Debug("FUNCEXIT")
}

func GetUser(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")
	var user models.User
	id := c.Param("id")
	database.Read(env.USER_COLLECTION, bson.M{"_id": id}).Decode(&user)
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	if sysAdmin || appUser.Role == "admin" || appUser.ID == id {
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to access this user details."})
	}
	Logger.Debug("FUNCEXIT")
}

func UpdateUser(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")
	var user models.User
	user.ID = c.Param("id")
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid json provided"})
	}
	userPassword := user.Password
	hashPassword, err := helpers.HashPassword(userPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	user.Password = hashPassword
	filter := bson.M{"_id": user.ID}
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	if sysAdmin || appUser.Role == "admin" {
		update := bson.M{"$set": bson.M{"first_name": user.First_Name, "last_name": user.Last_Name, "password": user.Password, "email": user.Email, "role": user.Role, "updated_at": time.Now().Local().String()}}
		response := database.Update(env.USER_COLLECTION, filter, update)
		c.JSON(http.StatusOK, gin.H{"Updated count": response.ModifiedCount, "Updated data": user})
	} else {
		if appUser.ID == user.ID {
			update := bson.M{"$set": bson.M{"first_name": user.First_Name, "last_name": user.Last_Name, "password": user.Password, "email": user.Email, "updated_at": time.Now().Local().String()}}
			response := database.Update(env.USER_COLLECTION, filter, update)
			c.JSON(http.StatusOK, gin.H{"Updated count": response.ModifiedCount, "Updated data": user})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to update user details."})
		}
	}
	Logger.Debug("FUNCEXIT")
}

func DeleteUser(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	role := c.GetString("role")
	username := c.GetString("username")
	password := c.GetString("password")
	id := c.Param("id")
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	if sysAdmin || appUser.Role == "admin" {
		result := database.Delete(env.USER_COLLECTION, bson.M{"_id": id})
		c.JSON(http.StatusOK, gin.H{"Deleted count": result.DeletedCount, "Deleted Count": result.DeletedCount})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to delete user"})
	}
	Logger.Debug("FUNCEXIT")
}
