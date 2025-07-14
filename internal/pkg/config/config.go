package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Log        Log        `mapstructure:"log"`
	AppConfig  AppConfig  `mapstructure:"app_config"`
	Connection Connection `mapstructure:"connection"`
}

type AppConfig struct {
	Env             string `mapstructure:"env"`
	Debug           bool   `mapstructure:"debug"`
	EnableProfiling bool   `mapstructure:"enable_profiling"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type Connection struct {
	HTTP       HTTP       `mapstructure:"http"`
	Postgresql Postgresql `mapstructure:"postgresql"`
}

type HTTP struct {
	Timeout int `mapstructure:"timeout"`
}

type Postgresql struct {
	Master PostgresqlInstance `mapstructure:"master"`
	Slave  PostgresqlInstance `mapstruture:"slave"`
}

type PostgresqlInstance struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	MaxOpen  int    `mapstructure:"max_open"`
	MaxIdle  int    `mapstructure:"max_idle"`
}

func NewConfig() (Config, error) {
	return newConfig([]string{".", "configs"})
}

func newConfig(paths []string) (Config, error) {
	conf := Config{}

	// Load .env files in each path if exists
	_ = godotenv.Load()
	for _, path := range paths {
		envPath := filepath.Join(path, ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
		}
	}

	// Create and configure viper
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Add config paths
	for _, path := range paths {
		v.AddConfigPath(path)
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("config.NewConfig: failed to read config file: %w", err)
	}

	// Unmarshal into config struct
	if err := v.Unmarshal(&conf); err != nil {
		return conf, fmt.Errorf("config.NewConfig: failed to unmarshal config: %w", err)
	}
	return conf, nil
}
