package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Environment        string `mapstructure:"ENVIRONMENT"`
	DBHost             string `mapstructure:"DB_HOST"`
	DBPort             string `mapstructure:"DB_PORT"`
	DBUser             string `mapstructure:"DB_USER"`
	DBName             string `mapstructure:"DB_NAME"`
	DBPassword         string `mapstructure:"DB_PASSWORD"`
	DBSSLMode          string `mapstructure:"DB_SSL_MODE"`
	JWTSecretKey       string `mapstructure:"JWT_SECRET_KEY"`
	AwsAccessKeyId     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AwsEndpoint        string `mapstructure:"AWS_ENDPOINT"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPort          string `mapstructure:"REDIS_PORT"`
	RedisPassword      string `mapstructure:"REDIS_PASSWORD"`
}

func bindEnv() {
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_SSL_MODE")
	viper.BindEnv("JWT_SECRET_KEY")
	viper.BindEnv("AWS_ACCESS_KEY_ID")
	viper.BindEnv("AWS_SECRET_ACCESS_KEY")
	viper.BindEnv("AWS_ENDPOINT")
	viper.BindEnv("REDIS_HOST")
	viper.BindEnv("REDIS_PORT")
	viper.BindEnv("REDIS_PASSWORD")
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
