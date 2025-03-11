package models

import (
    "time"
)

type User struct {
    ID           string    `json:"id,omitempty"             bson:"_id,omitempty"`
    Name         string    `json:"name"                     bson:"name,omitempty"           validate:"required"`
    Username     string    `json:"username"                 bson:"username,omitempty"       validate:"required"`
    Email        string    `json:"email"                    bson:"email,omitempty"          validate:"required,email"`
    Password     string    `json:"password"                 bson:"password,omitempty"       validate:"required,min=8"`
    CreatedAt    time.Time `json:"created_at"               bson:"created_at,omitempty"`
    ProfileImage string    `json:"profile_image,omitempty"  bson:"profile_image,omitempty"`
}
