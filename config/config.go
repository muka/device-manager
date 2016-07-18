package config

import (
	"log"
	"os"

	"github.com/muka/device-manager/util"

	"gopkg.in/yaml.v2"
)

const defaultPath string = "./config.yml"

// Config the parsed configuration object
type Config struct {
	APISpecPath string
}

var cfg *Config

func init() {
}

// Get returns the configuration struct
func Get() Config {
	if cfg == nil {
		Load(defaultPath)
	}
	return *cfg
}

// Load an yaml config file
func Load(path string) error {

	log.Printf("Load file %s", path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("File not found (%v)", err)
		return err
	}

	data, err := util.FromFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
		return err
	}

	cfg = new(Config)
	err = yaml.Unmarshal([]byte(data), cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
		return err
	}

	log.Printf("--- Loaded configuration:\n%v\n\n", cfg)

	return nil
}
