package api

import "github.com/PagerDuty/go-pagerduty"

type Team struct {
	ID   string
	Name string
}

func (p *PagerDutyClient) ListTeams() ([]*Team, error) {
	var opts pagerduty.ListTeamOptions
	listTeamsResponse, err := p.ApiClient.ListTeams(opts)
	if err != nil {
		return nil, err
	}

	var teamList []*Team
	for _, team := range listTeamsResponse.Teams {
		teamList = append(teamList, &Team{
			ID:   team.ID,
			Name: team.Name,
		})
	}
	return teamList, nil
}
