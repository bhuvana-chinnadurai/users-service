package conf

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel   string `mapstructure:"LOG_LEVEL"`
	ServerPort string `mapstructure:"SERVER_PORT"`

	DBUsername string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)

	switch os.Getenv("ENV") {
	case "PRODUCTION":
		viper.SetConfigName("production")
	case "DEVELOPMENT":
		viper.SetConfigName("development")
	case "TEST":
		viper.SetConfigName("test")
	}

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
