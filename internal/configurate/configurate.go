package configurate

import "os"

type Config struct {
	DBHost string
	DBPort string
	DBUSer string
	DBPass string
	DBName string
}

func New() Config {
	return Config{
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUSer: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
	}
}
