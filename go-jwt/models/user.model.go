package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	UserId       *string       `json:"user_id"`
	FirstName    *string       `json:"first_name" validate:"required, min=2, max=100"`
	LastName     *string       `json:"last_name" validate:"required, min=2, max=100"`
	PhoneNumber  *string       `json:"phone_number" validate:"required, min=10, max=10"`
	Email        *string       `json:"email" validate:"required, email"`
	Password     *string       `json:"password" validate:"required, min=6, max=100"`
	UserType     *string       `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Status       *string       `json:"status" validate:"eq=ACTIVE|eq=INACTIVE"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	DeletedAt    time.Time     `json:"deleted_at"`
	RefreshToken *string       `json:"refresh_token"`
	AccessToken  *string       `json:"access_token"` // best practice: do not store access token, but for practicing purpose, I will follow the code on YouTube
}
