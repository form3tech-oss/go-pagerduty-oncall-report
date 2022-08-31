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

type client interface {
	ListUsers() ([]*api.User, error)
	ListTeams() ([]*api.Team, error)
	ListServices(string) ([]*api.Service, error)
	ListSchedules() ([]*api.Schedule, error)
}

type pagerDutyClient struct {
	client client
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "configuration file (default is ~/.pd-report-config.yml)")

	viper.SetDefault("rotationStartHour", "08:00:00")
	viper.SetDefault("currency", "£")
}

func initConfig() {
	// Don't forget to read model either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use model file from the flag.
		viper.SetConfigFile(cfgFile)
		log.Println("Reading configuration file:", cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal("Can't get the homedir: ", err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".pd-report-config")
		log.Println("Reading configuration file:", fmt.Sprintf("%s/.pd-report-config-yml", home))
	}

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Can't read config: ", err)
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
