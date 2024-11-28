package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName  string         `mapstructure:"app_name"`
	Debug    bool           `mapstructure:"debug"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"DB"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	// Load YAML config
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	viper.AutomaticEnv()

	viper.SetConfigFile(".env-example")
	viper.SetConfigType("env")
	if err := viper.MergeInConfig(); err != nil {
		fmt.Printf("Warning: .env-example file not found: %v\n", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %w", err)
	}

	fmt.Printf("Loaded Config: %+v\n", config)
	return config, nil
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}
