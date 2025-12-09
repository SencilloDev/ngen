package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfg Config

var rootCmd = &cobra.Command{
	Use:   "ngen",
	Short: "NATS micro service diagramming and OpenAPI spec creation",
}
var replacer = strings.NewReplacer("-", "_")

type Config struct {
	ServiceName  string `mapstructure:"service_name"`
	ServiceID    string `mapstructure:"service_id"`
	MethodOffset int    `mapstructure:"method_offset"`
	NatsURLs     string `mapstructure:"nats_urls"`
	NatsSeed     string `mapstructure:"nats_seed"`
	NatsJWT      string `mapstructure:"nats_jwt"`
	NatsSecret   string `mapstructure:"nats_secret"`
	CredsFile    string `mapstructure:"credentials_file"`
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.natsoapi.json)")
}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".natsoapi")
	}

	viper.SetEnvPrefix("natsoapi")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := viper.ReadInConfig(); err == nil {
		logger.Debug(fmt.Sprintf("using config %s", viper.ConfigFileUsed()))
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		cobra.CheckErr(err)
	}
}
