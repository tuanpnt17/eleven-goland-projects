package helpers

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tuanpnt17/eleven-golang-projects/go-jwt/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"os"
	"time"
)

type SignedDetails struct {
	Uid       string
	Email     string
	FirstName string
	LastName  string
	UserType  string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var SecretKey string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstname, lastname, userType, uid string) (accessToken, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		UserType:  userType,
		Uid:       uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(168))),
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SecretKey))

	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func UpdateAllToken(accessToken, refreshToken, uid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj bson.D

	updateObj = append(updateObj, bson.E{Key: "access_token", Value: accessToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: refreshToken})

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updateAt})

	filter := bson.M{"user_id": uid}
	opts := options.UpdateOne().SetUpsert(true)

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		opts,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
	}
	return
}
