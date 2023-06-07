package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	CLICKHOUSE_DB_USERNAME string `mapstructure:"CLICKHOUSE_DB_USERNAME"`
	CLICKHOUSE_DB_PASSWORD string `mapstructure:"CLICKHOUSE_DB_PASSWORD"`
	CLICKHOUSE_DB_HOST     string `mapstructure:"CLICKHOUSE_DB_HOST"`
	CLICKHOUSE_DB_NAME     string `mapstructure:"CLICKHOUSE_DB_NAME"`
	MIGRATE_PATH           string `mapstructure:"MIGRATE_PATH"`

	SERVER_HOST string `mapstructure:"SERVER_HOST"`

	GRPC_USER_SERVICE_HOST   string `mapstructure:"GRPC_USER_SERVICE_HOST"`
	GRPC_DRIVER_SERVICE_HOST string `mapstructure:"GRPC_DRIVER_SERVICE_HOST"`
	GRPC_ORDER_SERVICE_HOST  string `mapstructure:"GRPC_ORDER_SERVICE_HOST"`

	ADMIN_LOGIN string `mapstructure:"ADMIN_LOGIN"`
	ADMIN_PASS  string `mapstructure:"ADMIN_PASS"`

	HS256_SECRET string `mapstructure:"HS256_SECRET"`

	BROKER_HOST string `mapstructure:"BROKER_HOST"`

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

func (c *Config) GetClickhouseUrl() string {
	return fmt.Sprintf("clickhouse://%s?username=%s&password=%s&database=%s&x-multi-statement=true", c.CLICKHOUSE_DB_HOST, c.CLICKHOUSE_DB_USERNAME, c.CLICKHOUSE_DB_PASSWORD, c.CLICKHOUSE_DB_NAME)
}
