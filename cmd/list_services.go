package cmd

import (
	"fmt"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"

	"github.com/spf13/cobra"
)

var listServicesCmd = &cobra.Command{
	Use:   "services <teamID>",
	Short: "list services on PagerDuty",
	Long:  "Get the list of services configured in PagerDuty",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pd := &pagerDutyClient{client: api.NewPagerDutyAPIClient(Config.PdAuthToken)}
		return pd.listServices(args[0])
	},
}

func init() {
	rootCmd.AddCommand(listServicesCmd)
}

func (pd *pagerDutyClient) listServices(teamID string) error {
	services, err := pd.client.ListServices(teamID)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("==== Found %d service(s) for the team %s ====", len(services), teamID))
	for _, service := range services {
		fmt.Println(fmt.Sprintf("[%s] %-20s", service.ID, service.Name))
	}

	return nil
}
