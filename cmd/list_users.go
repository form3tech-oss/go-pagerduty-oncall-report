package cmd

import (
	"fmt"

	"github.com/rogersole/go-pagerduty-oncall-report/api"
	"github.com/spf13/cobra"
)

var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "list users on PagerDuty",
	Long:  "Get the list of users configured in PagerDuty",
	RunE:  listUsers,
}

func init() {
	rootCmd.AddCommand(listUsersCmd)
}

func listUsers(cmd *cobra.Command, args []string) error {
	users, err := api.Client.ListUsers()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("==== Found %d user(s) ====", len(users)))
	for _, user := range users {
		var userTeams string
		for _, userTeam := range user.Teams {
			userTeams += fmt.Sprintf("%s ", userTeam.ID)
		}
		fmt.Println(fmt.Sprintf("[%s] %-20s %-30s in teams: %s", user.ID, user.Name, fmt.Sprintf("<%s>", user.Email), userTeams))
	}

	return nil
}
