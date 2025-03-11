package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Recipe struct {
    ID           primitive.ObjectID `json:"id,omitempty"    bson:"_id,omitempty" `
    Title        string             `json:"title"           bson:"title,omitempty"          validate:"required"`
    Image        []byte             `json:"image_url"       bson:"image_url,omitempty"      validate:"required"`
    Ingredients  []string           `json:"ingredients"     bson:"ingredients,omitempty"    validate:"required"`
    Instructions []string           `json:"instructions"    bson:"instructions,omitempty"   validate:"required"`
    PrepTime     string             `json:"prep_time"       bson:"prep_time,omitempty"      validate:"required,min=1"`
    Servings     string             `json:"servings"        bson:"servings,omitempty"       validate:"required,min=1"`
    Difficulty   string             `json:"difficulty"      bson:"difficulty,omitempty"     validate:"required"`
    CreatedAt    time.Time          `json:"created_at"      bson:"created_at,omitempty"`
    CreatedBy    string             `json:"created_by"      bson:"created_by,omitempty"`
}
