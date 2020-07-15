package api

import (
	"io/ioutil"
	"log"

	"github.com/rasha108/apiCargoRest.git/internal/app/rabbitmq"
	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLeval   string          `yaml:"log_level"`
	SessionKey string          `yaml:"session_key"`
	APIConfig  APIConfig       `yaml:"api"`
	MailConfig rabbitmq.Config `yaml:"mail"`
	DbConfig   DbConfig        `yaml:"db"`
}

// APIConfig stores api parameters
type APIConfig struct {
	Bind     string `yaml:"bind"`
	BasePath string `yaml:"basePath"`
}

// DbConfig stores database connection parametrs
type DbConfig struct {
	Host           string `yaml:"host"`
	User           string `yaml:"user"`
	Pass           string `yaml:"pass"`
	DbName         string `yaml:"dbName"`
	Port           int    `yaml:"port"`
	MaxConnections int    `yaml:"maxConnections"`
}

//GetCOnfig initialize configuration by path
func GetConfig(filepath string) (*Config, error) {
	var config Config
	config.MailConfig.SetDefaults()

	configContent, err := ioutil.ReadFile(filepath)
	if err == nil {
		err = yaml.Unmarshal(configContent, &config)
		if err != nil {
			log.Fatalf("unmarshal '%s' content failed: '%s'", filepath, err)
		}
	} else {
		log.Fatalf("read config '%s' failed: '%s'", filepath, err)
	}
	return &config, err
}
