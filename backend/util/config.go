package util

import "github.com/spf13/viper"

type Config struct {
	ServerHost                  string `mapstructure:"server_host"`
	ServerPort                  string `mapstructure:"server_port"`
	DbUrl                       string `mapstructure:"db_url"`
	FrontOrigin                 string `mapstructure:"front_origin"`
	RequestContentLengthMaxByte int    `mapstructure:"request_content_length_max_byte"`
}

func LoadConfig() (cfg Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return
	}
	return
}
