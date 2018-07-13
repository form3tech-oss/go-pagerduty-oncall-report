package cmd

import (
	"fmt"

	"github.com/rogersole/go-pagerduty-oncall-report/api"
	"github.com/spf13/cobra"
)

var listTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "list teams on PagerDuty",
	Long:  "Get the list of teams configured in PagerDuty",
	RunE:  listTeams,
}

func init() {
	rootCmd.AddCommand(listTeamsCmd)
}

func listTeams(cmd *cobra.Command, args []string) error {
	teams, err := api.Client.ListTeams()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("==== Found %d team(s) ====", len(teams)))
	for _, team := range teams {
		fmt.Println(fmt.Sprintf("[%s] %-20s", team.ID, team.Name))
	}

	return nil
}
