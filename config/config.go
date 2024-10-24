package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DbPath string
}

func Load() *Config {
	var config Config

	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка при чтении файла конфигурации: %v", err)
	}

	var dbPath = viper.Get("TEST_DB_PATH")
	if dbPath == "" {
		panic("TEST_DB_PATH is empty")
	}

	config.DbPath = dbPath.(string)
	return &config
}
