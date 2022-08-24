package cmd

import (
	"fmt"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"

	"github.com/spf13/cobra"
)

type client interface {
	ListUsers() ([]*api.User, error)
}

type pagerDutyClient struct {
	client client
}

var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List users on PagerDuty",
	Long:  "Get the list of users configured in PagerDuty",
	RunE: func(cmd *cobra.Command, args []string) error {
		pd := &pagerDutyClient{client: api.NewPagerDutyAPIClient(Config.PdAuthToken)}
		err := pd.listUsers()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listUsersCmd)
}

func (pd *pagerDutyClient) listUsers() error {
	users, err := pd.client.ListUsers()
	if err != nil {
		return fmt.Errorf("failed to fetch user list: %w", err)
	}

	fmt.Println(fmt.Sprintf("==== Found %d user(s) ====", len(users)))
	for _, user := range users {
		var userTeams string
		for _, userTeam := range user.Teams {
			userTeams += fmt.Sprintf("%s ", userTeam.ID)
		}
		output := fmt.Sprintf("[%s] %-30s %-40s in teams: %s", user.ID, user.Name, fmt.Sprintf("<%s>", user.Email), userTeams)
		fmt.Println(output)
	}

	return nil
}
