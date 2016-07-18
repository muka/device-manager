package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// MaxMemory to be read at once from an http endpoint
const MaxMemory = 1024

func check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

// FromURL load the content of an url
func FromURL(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	buffer, err := ioutil.ReadAll(io.LimitReader(response.Body, MaxMemory))
	return string(buffer), err
}

// FromFile load a file as text
func FromFile(src string) (string, error) {
	res, err := ioutil.ReadFile(src)
	check(err)
	fmt.Print(string(res))
	return string(res), nil
}
