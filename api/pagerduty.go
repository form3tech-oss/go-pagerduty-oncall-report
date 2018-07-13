package api

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

var Client *PagerDutyClient

type PagerDutyClient struct {
	apiClient *pagerduty.Client
}

type ScheduleInfo struct {
	ID            string
	Name          string
	Location      *time.Location
	Start         time.Time
	End           time.Time
	FinalSchedule pagerduty.ScheduleLayer
}

type UserRotaPeriod struct {
	Start time.Time
	End   time.Time
}

type UserRotaInfo struct {
	ID      string
	Name    string
	Periods []*UserRotaPeriod
}

type ScheduleUserRotationData map[string]*UserRotaInfo

func InitialisePagerDutyAPIClient(authToken string) {
	Client = &PagerDutyClient{
		apiClient: pagerduty.NewClient(authToken),
	}
}

func (p *PagerDutyClient) ListSchedules() ([]pagerduty.Schedule, error) {
	var opts pagerduty.ListSchedulesOptions
	listSchedulesResponse, err := p.apiClient.ListSchedules(opts)
	if err != nil {
		return nil, err
	}

	return listSchedulesResponse.Schedules, nil
}

func (p *PagerDutyClient) ListServices(teamID string) ([]pagerduty.Service, error) {
	var opts pagerduty.ListServiceOptions
	opts.TeamIDs = []string{teamID}
	listServicesResponse, err := p.apiClient.ListServices(opts)
	if err != nil {
		return nil, err
	}

	return listServicesResponse.Services, nil
}

func (p *PagerDutyClient) ListTeams() ([]pagerduty.Team, error) {
	var opts pagerduty.ListTeamOptions
	listTeamsResponse, err := p.apiClient.ListTeams(opts)
	if err != nil {
		return nil, err
	}
	return listTeamsResponse.Teams, nil
}

func (p *PagerDutyClient) ListUsers() ([]pagerduty.User, error) {
	var opts pagerduty.ListUsersOptions
	listUsersResponse, err := p.apiClient.ListUsers(opts)
	if err != nil {
		return nil, err
	}
	return listUsersResponse.Users, nil
}

func (p *PagerDutyClient) GetSchedule(scheduleID, startDate, endDate string) (*pagerduty.Schedule, error) {
	var opts pagerduty.GetScheduleOptions
	opts.Since = startDate
	opts.Until = endDate
	scheduleResponse, err := p.apiClient.GetSchedule(scheduleID, opts)
	if err != nil {
		return nil, err
	}
	return scheduleResponse, nil
}
