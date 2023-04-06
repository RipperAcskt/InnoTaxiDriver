package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	CASSANDRA_DB_USERNAME string `mapstructure:"CASSANDRA_DB_USERNAME"`
	CASSANDRA_DB_PASSWORD string `mapstructure:"CASSANDRA_DB_PASSWORD"`
	CASSANDRA_CLUSTER     string `mapstructure:"CASSANDRA_CLUSTER"`
	CASSANDRA_KEYSPACE    string `mapstructure:"CASSANDRA_KEYSPACE"`

	MIGRATE_PATH string `mapstructure:"MIGRATE_PATH"`

	SALT                   string `mapstructure:"SALT"`
	GRPC_USER_SERVICE_HOST string `mapstructure:"GRPC_USER_SERVICE_HOST"`

	HS256_SECRET string `mapstructure:"HS256_SECRET"`

	GRPC_HOST string `mapstructure:"GRPC_HOST"`
}

func New() (*Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}
	return config, nil
}
