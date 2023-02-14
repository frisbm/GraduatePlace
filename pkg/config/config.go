package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBName       string `mapstructure:"DB_NAME"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBSSLMode    string `mapstructure:"DB_SSL_MODE"`
	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
}

func bindEnv() {
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_SSL_MODE")
	viper.BindEnv("JWT_SECRET_KEY")
}

func NewConfig(filepath string) *Config {
	config := new(Config)
	bindEnv()
	v := viper.New()
	v.SetConfigFile(filepath)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Default().Printf("Using config file: %s", filepath)
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return config
}
