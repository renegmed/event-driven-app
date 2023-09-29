package config

import "github.com/spf13/viper"

type SubscriberConfig struct {
	Group string
	Queue string
}

type Config struct {
	Jobs SubscriberConfig
}

func LoadConfigFrom(in string) (Config, error) {
	conf := viper.New()
	conf.SetConfigName("config")
	conf.SetConfigType("yml")
	conf.AddConfigPath(in)

	if err := conf.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	err := conf.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func LoadConfig() (Config, error) {
	return LoadConfigFrom(".")
}
