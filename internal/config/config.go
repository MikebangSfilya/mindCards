package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	HTTTPServer `yaml:"http_server"`
	ConfigDB    `yaml:"DB_CFG"`
}

type HTTTPServer struct {
	Adress      string        `yaml:"addres" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"30s"`
}

type ConfigDB struct {
	DBHost string `yaml:"DB_HOST"`
	DBPort string `yaml:"DB_PORT"`
	DBUSer string `yaml:"DB_USER"`
	DBPass string `yaml:"DB_PASSWORD"`
	DBName string `yaml:"DB_NAME"`
}

func NewDB() ConfigDB {
	return ConfigDB{
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUSer: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
	}
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file not exitst")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return cfg
}

func (c ConfigDB) Get() []string {
	envParams := make([]string, 0)
	envParams = append(envParams, c.DBHost, c.DBPort, c.DBUSer, c.DBPass, c.DBName)
	return envParams
}
