package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	LoggerLevel    string `mapstructure:"DAAS_LOGGER_LEVEL"`
	LoggerEncoding string `mapstructure:"DAAS_LOGGER_ENCODING"`
	APIAddress     string `mapstructure:"DAAS_API_ADDRESS"`

	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
}

var vp *viper.Viper

func GetConfig() (*Config, error) {
	vp = viper.New()
	var config Config

	vp.SetConfigName("config")
	vp.SetConfigType("env")
	vp.AddConfigPath(".")

	err := vp.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return &Config{}, err
	}

	// Read in environment variables
	vp.AutomaticEnv()

	err = vp.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct: ", err)
		return &Config{}, err
	}
	return &config, nil
}
