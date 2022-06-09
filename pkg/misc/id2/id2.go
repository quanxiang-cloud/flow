package id2

import (
	"strings"

	"github.com/google/uuid"
)

// ID generator

// GenID gen id
func GenID() string {
	return gUUID()
}

// GenUpperID gen upper id
func GenUpperID() string {
	return strings.ToUpper(
		strings.Replace(
			gUUID(),
			"-", "", -1))
}

// gUUID google/uuid
func gUUID() string {
	return uuid.New().String()
}
