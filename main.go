package main

import (
	"fmt"
	"log"

	"github.com/muka/device-manager/config"
)

func main() {

	path := "./config.yaml"
	fmt.Printf("Loading configuration (src: %s)\n", path)
	conf, err := config.Load(path)
	if err != nil {
		log.Fatalf("Error loading configuration: %v\n", err)
		panic(err)
	}
	fmt.Printf("Config %v\n", conf)

}
