package config

import "github.com/spf13/viper"

func GetStringValue(key string, defaultValue string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	return defaultValue
}

func GetIntValue(key string, defaultValue int) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}

	return defaultValue
}
