package pkg

import (
	"github.com/spf13/viper"
	"log"
)

func EnvVar(key string) string {
	viper.SetConfigFile("app.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	return value
}
