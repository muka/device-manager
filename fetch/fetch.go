package fetch

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/muka/device-manager/util"
)

var logger = util.Logger()

// MaxMemory to be read at once from an http endpoint
const MaxMemory = 1024

func check(err error) {
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
}

// FromURL load the content of an url
func FromURL(url string) (string, error) {
	logger.Printf("Load URL %s", url)
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
	logger.Printf("Load file %s", src)
	res, err := ioutil.ReadFile(src)
	check(err)
	return string(res), nil
}

// GetContent load content from an URL or file
func GetContent(src string) (string, error) {
	if strings.Contains(src, "http") {
		return FromURL(src)
	}
	return FromFile(src)
}
