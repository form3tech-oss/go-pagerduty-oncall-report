package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"
	"github.com/form3tech-oss/go-pagerduty-oncall-report/configuration"
)

var (
	cfgFile string
	Config  *configuration.Configuration
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "configuration file (default is ~/.pd-report-config.yml)")

	viper.SetDefault("rotationStartHour", "08:00:00")
	viper.SetDefault("currency", "£")
}

func initConfig() {
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal("failed to get home directory: ", err)
		}

		filename := ".pd-report-config.yml"
		cfgFile = fmt.Sprintf("%s/%s", home, filename)
	}

	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")
	log.Println("reading configuration file:", cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to read config file: ", err)
	}

	Config = configuration.New()
	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("%v, %#v", err, Config)
	}

	api.InitialisePagerDutyAPIClient(Config.PdAuthToken)
}

var rootCmd = &cobra.Command{
	Use:   "pd-report",
	Short: "Easily generate PagerDuty reports",
	Long: `Generate on-call rotation reports automatically
from your PagerDuty account.`,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
