package cmd

import (
	"fmt"

	"github.com/rogersole/go-pagerduty-oncall-report/api"
	"github.com/spf13/cobra"
)

var listServicesCmd = &cobra.Command{
	Use:   "services <teamID>",
	Short: "list services on PagerDuty",
	Long:  "Get the list of services configured in PagerDuty",
	Args:  cobra.ExactArgs(1),
	RunE:  listServices,
}

func init() {
	rootCmd.AddCommand(listServicesCmd)
}

func listServices(cmd *cobra.Command, args []string) error {
	teamID := args[0]
	services, err := api.Client.ListServices(teamID)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("==== Found %d service(s) for the team %s ====", len(services), teamID))
	for _, service := range services {
		fmt.Println(fmt.Sprintf("[%s] %-20s", service.ID, service.Name))
	}

	return nil
}
