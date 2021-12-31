package env

import (
	"service-discovery/middlewares"

	"github.com/spf13/viper"
)

var Logger = middlewares.Logger()

func GetEnvironmentVariable(key string) string {
	viper.SetConfigFile(".env")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		Logger.Error("Error while reading config file " + err.Error())
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		Logger.Error("Invalid type assertion")
	}
	return value
}
