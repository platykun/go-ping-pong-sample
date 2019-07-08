package service

import (
	"github.com/google/uuid"
	"strings"
)

// GenerateID returns unique id
func GenerateID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
