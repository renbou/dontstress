package utils

import "github.com/google/uuid"

func GetId() string {
	return uuid.New().String()
}
