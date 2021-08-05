package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DbDriver            string        `mapstructure:"DB_DRIVER"`
	DbSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenAccessDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
