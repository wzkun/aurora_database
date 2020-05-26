package utils

import (
	uuid "github.com/satori/go.uuid"
)

// NewUUIDV4 return uuid string of version 4
func NewUUIDV4() string {
	uid := uuid.NewV4()

	return uid.String()
}
