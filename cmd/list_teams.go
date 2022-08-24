package cmd

import (
	"fmt"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"

	"github.com/spf13/cobra"
)

var listTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "list teams on PagerDuty",
	Long:  "Get the list of teams configured in PagerDuty",
	RunE: func(cmd *cobra.Command, args []string) error {
		pd := &pagerDutyClient{client: api.NewPagerDutyAPIClient(Config.PdAuthToken)}
		err := pd.listTeams()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listTeamsCmd)
}

func (pd *pagerDutyClient) listTeams() error {
	teams, err := pd.client.ListTeams()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("==== Found %d team(s) ====", len(teams)))
	for _, team := range teams {
		fmt.Println(fmt.Sprintf("[%s] %-20s", team.ID, team.Name))
	}

	return nil
}
