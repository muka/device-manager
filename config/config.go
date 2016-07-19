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
	APISpecPath string `yaml:"XMLSpecPath"`
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
		logger.Fatalf("File not found (%v)", err)
		return err
	}

	data, err := fetch.GetContent(path)
	if err != nil {
		logger.Fatalf("error: %v", err)
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
