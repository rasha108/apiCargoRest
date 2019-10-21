package api

import (
	"io/ioutil"
	"log"

	"github.com/rasha108/apiCargoRest.git/internal/app/rabbitclient"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BindAddr    string              `yaml:"bind_addr"`
	LogLeval    string              `yaml:"log_level"`
	DatabaseURL string              `yaml:"database_url"`
	SessionKey  string              `yaml:"session_key"`
	BasePath    string              `yaml:"base_path"`
	MailConfig  rabbitclient.Config `yaml:"mail"`
	DbConfig    DbConfig            `yaml:"db"`
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
