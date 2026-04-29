package models

import (
	"time"
	"github.com/google/uuid"
)

//the user body
type User struct {
    ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	First_Name    *string   `json:"first_name" validate:"required,min=3,max=20"`
	Last_Name     *string   `json:"last_name" validate:"required,min=2,max=10"`
	Email         *string   `json:"email" validate:"required,email"`
	Password      *string   `json:"password" validate:"required,min=4"`
	Token         *string   `json:"token,omitempty"`
	Refresh_token *string   `json:"refresh_token,omitempty"`
	User_type     *string   `json:"user_type" validate:"required,oneof=ADMIN USER"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
   //User_id       string    `json:"user_id"`
}
