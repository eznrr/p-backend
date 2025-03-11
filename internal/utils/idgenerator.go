package utils

import "github.com/google/uuid"

func GenerateStringID() string {
	id, err := uuid.NewV7()

	if err != nil {
		id = uuid.New()
	}

	return id.String()
}

func ValidateID(id string) bool {
	err := uuid.Validate(id)
	return err == nil 
}
