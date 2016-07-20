package util

import "log"

// CheckError panic if err is not nil
func CheckError(err error) {
	if err != nil {
		log.Fatalf("Error occured: \n\n%v\n\n", err)
		panic(err)
	}
}
