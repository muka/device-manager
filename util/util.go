package util

import (
	"github.com/satori/go.uuid"
	"log"
	"strings"
)

// CheckError panic if err is not nil
func CheckError(err error) {
	if err != nil {
		log.Fatalf("An error occured: \n%v\n\n", err)
		panic(err)
	}
}

// GenerateID generate an unique id based on UUIDv4
func GenerateID() string {
	uuid := uuid.NewV4().String()
	return strings.Replace(uuid, "-", "_", -1)
}
