package config

import (
	"fmt"
	"log"

	"github.com/muka/device-manager/util"

	"gopkg.in/yaml.v2"
)

// T the parsed configuration object
type T struct {
	APISpecPath string
}

// Load an yaml config file
func Load(path string) (t T, err error) {

	data, err := util.FromFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	fmt.Printf("--- t:\n%v\n\n", t)

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	fmt.Printf("--- t dump:\n%s\n\n", string(d))

	return t, nil
}
