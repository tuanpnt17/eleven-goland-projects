package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tuanpnt17/eleven-golang-projects/go-jwt/database"
	"github.com/tuanpnt17/eleven-golang-projects/go-jwt/helpers"
	"github.com/tuanpnt17/eleven-golang-projects/go-jwt/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var userCollection = database.OpenCollection(database.Client, "users")
var validate = validator.New(validator.WithRequiredStructEnabled())

func HashPassword(userPassword string) string {
	password, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}
	return string(password)
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		}

		accessToken, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, *foundUser.UserId)
		helpers.UpdateAllToken(accessToken, refreshToken, *foundUser.UserId)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.UserId}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, foundUser)

	}
}

func VerifyPassword(inputPassword string, foundPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(inputPassword))
	check := true
	msg := ""
	if err != nil {
		check = false
		msg = fmt.Sprintf("Email or Password is incorrect")
	}
	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Panic(err)
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.PhoneNumber})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Panic(err)
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Phone Number already exists"})
			return
		}

		// create user here
		// Todo: Need to separate it to a user service

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) // RFC3339 is just the time format
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = bson.NewObjectID()

		userId := user.ID.Hex()
		user.UserId = &userId

		accessToken, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *user.UserId)
		user.AccessToken = &accessToken
		user.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created. Error: %s", insertErr.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

const (
	timeoutDuration   = 100 * time.Second
	defaultPageSize   = 10
	defaultStartIndex = 0
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create a context with timeout
		timeoutCtx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		// Fetch pagination parameters
		startIndex, recordsPerPage := getPaginationParams(c)

		// Build aggregation pipeline for MongoDB
		aggregationPipeline := buildAggregationPipeline(startIndex, recordsPerPage)

		// Fetch aggregate data from the database
		aggregate, err := userCollection.Aggregate(timeoutCtx, aggregationPipeline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse results
		var users []bson.M
		if err = aggregate.All(timeoutCtx, &users); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Send the response
		c.JSON(http.StatusOK, users[0])
	}
}

func getPaginationParams(c *gin.Context) (int, int) {
	recordsPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordsPerPage < 1 {
		recordsPerPage = defaultPageSize
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex, err := strconv.Atoi(c.Query("startIndex"))
	if err != nil || startIndex < 0 {
		startIndex = (page - 1) * recordsPerPage
	}

	return startIndex, recordsPerPage
}

func buildAggregationPipeline(startIndex, recordsPerPage int) bson.A {
	matchStage := bson.D{{"$match", bson.D{{}}}}
	groupStage := bson.D{{"$group", bson.D{
		{"_id", "null"},
		{"total_count", bson.D{{"$sum", 1}}},
		{"data", bson.D{{"$push", "$$ROOT"}}},
	}}}
	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordsPerPage}}}},
		}},
	}
	return bson.A{matchStage, groupStage, projectStage}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
