package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName    string
	GitVersion     string
	BaseURL        string
	Port           string
	RunMode        string
	AllowedOrigins string
	DBHost         string
	DBPort         string
	DBUser         string
	DBName         string
	DBPass         string
	RedisHost      string
	RedisPort      string
	RedisPass      string
}

// GetCORS in config
func (c *Config) GetCORS() []string {
	cors := strings.Split(c.AllowedOrigins, ";")
	rs := []string{}
	for idx := range cors {
		itm := cors[idx]
		if strings.TrimSpace(itm) != "" {
			rs = append(rs, itm)
		}
	}
	return rs
}

// Loader load config from reader into Viper
type Loader interface {
	Load(viper.Viper) (*viper.Viper, error)
}

// generateConfigFromViper generate config from viper data
func generateConfigFromViper(v *viper.Viper) Config {
	return Config{
		ServiceName: v.GetString("SERVICE_NAME"),
		Port:        v.GetString("PORT"),
		BaseURL:     v.GetString("BASE_URL"),
		RunMode:     v.GetString("RUN_MODE"),
		GitVersion:  v.GetString("GIT_VERSION"),

		AllowedOrigins: v.GetString("ALLOWED_ORIGINS"),

		DBHost: v.GetString("DB_HOST"),
		DBPort: v.GetString("DB_PORT"),
		DBUser: v.GetString("DB_USER"),
		DBName: v.GetString("DB_NAME"),
		DBPass: v.GetString("DB_PASS"),

		RedisHost: v.GetString("REDIS_HOST"),
		RedisPort: v.GetString("REDIS_PORT"),
		RedisPass: v.GetString("REDIS_PASS"),
	}
}

// DefaultConfigLoaders is default loader list
func DefaultConfigLoaders() []Loader {
	loaders := []Loader{}
	fileLoader := NewFileLoader(".env", ".")
	loaders = append(loaders, fileLoader)
	loaders = append(loaders, NewENVLoader())

	return loaders
}

// LoadConfig load config from loader list
func LoadConfig(loaders []Loader) Config {
	v := viper.New()
	v.SetDefault("PORT", "8080")
	v.SetDefault("RUN_MODE", "local")

	for idx := range loaders {
		newV, err := loaders[idx].Load(*v)

		if err == nil {
			v = newV
		}
	}
	return generateConfigFromViper(v)
}
