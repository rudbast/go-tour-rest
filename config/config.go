package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	GO_ENV string = "GO_ENV"
)

var (
	ConfigFilePath = map[string]string{
		"production":  "config/production.ini",
		"development": "config/development.ini",
	}

	Data Config
)

type Config struct {
	JWT      JwtConfig
	Database DatabaseConfig
}

type JwtConfig struct {
	PublicKeyPath  string
	PrivateKeyPath string
}

type DatabaseConfig struct {
	Driver           string
	ConnectionString string
}

// Initialize configurations.
func Init() error {
	var env string

	if env = os.Getenv(GO_ENV); env == "" {
		env = "development"
	}

	if err := LoadConfig(env); err != nil {
		return err
	}

	return nil
}

// Load configuration from file based on given environment setting.
func LoadConfig(env string) error {
	configFile, err := ioutil.ReadFile(ConfigFilePath[env])

	if err != nil {
		log.Panic(err)
		return err
	}

	if _, err := toml.Decode(string(configFile), &Data); err != nil {
		log.Panic(err)
		return err
	}

	fmt.Println("Config loaded.")
	return nil
}
