package api

import "github.com/PagerDuty/go-pagerduty"

type Team struct {
	ID          string
	Name        string
	Description string
}

func (p *PagerDutyClient) ListTeams() ([]Team, error) {
	var opts pagerduty.ListTeamOptions
	listTeamsResponse, err := p.ApiClient.ListTeams(opts)
	if err != nil {
		return nil, err
	}

	var teams []Team

	for _, team := range listTeamsResponse.Teams {
		teams = append(teams, Team{
			ID:          team.ID,
			Name:        team.Name,
			Description: team.Description,
		})
	}
	return teams, nil
}
