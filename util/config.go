package util

import "github.com/spf13/viper"

// Config struct will hold all the conf vars that we will read from file or env variables
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig will read configurations from a config file if exists inside the path
// or override their values with env variables if provided
func LoadConfig(path string) (config Config, err error)  {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig(); if err != nil{
		return
	}
	err = viper.Unmarshal(&config)
	return
}