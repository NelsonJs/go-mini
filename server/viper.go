package server

import "github.com/spf13/viper"

func init() {
	viper.SetConfigFile("./config.toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
