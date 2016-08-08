package config

import (
	"os"

	"github.com/muka/device-manager/fetch"
	"github.com/muka/device-manager/util"

	"gopkg.in/yaml.v2"
)

const defaultPath string = "./config.yml"

// Config the parsed configuration object
type Config struct {
	debug bool              `yaml:"debug"`
	paths map[string]string `yaml:"paths"`
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

	logger := util.Logger()

	logger.Println("Load configuration")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Fatalf("Configuration file not found (%v)\n", err)
		return err
	}

	data, err := fetch.GetContent(path)
	if err != nil {
		logger.Fatalf("Cannot read configuration file: \n%v\n", err)
		return err
	}

	cfg = new(Config)
	err = yaml.Unmarshal([]byte(data), cfg)
	if err != nil {
		logger.Fatalf("error: %v", err)
		return err
	}

	logger.Printf("Configuration loaded")

	return nil
}
