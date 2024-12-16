package config

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type BotConfig struct {
	Token   string
	Verbose bool
}

type ServerConfig struct {
	ListenPort string `mapstructure:"listen_port"`
	Endpoint   string
}

type JWTConfig struct {
	Secret string
}

type Config struct {
	Bot    BotConfig
	Server ServerConfig
	Debug  bool
	JWT    JWTConfig
}

var GlobalConfig *Config

func LoadConfig(cmd *cobra.Command) error {
	GlobalConfig = &Config{}

	viper.SetDefault("debug", false)
	viper.SetDefault("bot.token", "")
	viper.SetDefault("bot.endpoint", "")
	viper.SetDefault("bot.verbose", false)
	viper.SetDefault("bot.listen_port", "8002")
	viper.SetDefault("db.uri", "")
	viper.SetDefault("db.name", "")
	viper.SetDefault("server.listen_port", "8000")
	viper.SetDefault("server.endpoint", "localhost:8000")
	viper.SetDefault("jwt.secret", "")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return err
	}

	return nil
}
