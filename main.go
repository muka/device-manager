package main

import (
	"fmt"

	"github.com/muka/device-manager/config"
)

func main() {

	config := config.Get()
	fmt.Printf("Config %v\n", config)

}
