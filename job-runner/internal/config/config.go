package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Host string
	Port string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Config struct {
	Server   ServerConfig
	RabbitMQ RabbitMQConfig
	DB       DBConfig
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
