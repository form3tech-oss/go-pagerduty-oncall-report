package api

import "github.com/PagerDuty/go-pagerduty"

type Service struct {
	ID   string
	Name string
}

func (p *PagerDutyClient) ListServices(teamID string) ([]*Service, error) {
	var opts pagerduty.ListServiceOptions
	opts.TeamIDs = []string{teamID}
	listServicesResponse, err := p.ApiClient.ListServices(opts)
	if err != nil {
		return nil, err
	}

	var serviceList []*Service
	for _, service := range listServicesResponse.Services {
		serviceList = append(serviceList, &Service{
			ID:   service.ID,
			Name: service.Name,
		})
	}

	return serviceList, nil
}
