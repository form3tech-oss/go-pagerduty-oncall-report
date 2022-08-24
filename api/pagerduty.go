package api

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

var Client *PagerDutyClient

type PdClient interface {
	ListSchedules(o pagerduty.ListSchedulesOptions) (*pagerduty.ListSchedulesResponse, error)
	ListServices(o pagerduty.ListServiceOptions) (*pagerduty.ListServiceResponse, error)
	ListTeams(o pagerduty.ListTeamOptions) (*pagerduty.ListTeamResponse, error)
	ListUsers(o pagerduty.ListUsersOptions) (*pagerduty.ListUsersResponse, error)
	GetUser(id string, o pagerduty.GetUserOptions) (*pagerduty.User, error)
	GetSchedule(id string, o pagerduty.GetScheduleOptions) (*pagerduty.Schedule, error)
}

type PagerDutyClient struct {
	ApiClient PdClient
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
		ApiClient: pagerduty.NewClient(authToken),
	}
}

func NewPagerDutyAPIClient(authToken string) *PagerDutyClient {
	return &PagerDutyClient{
		ApiClient: pagerduty.NewClient(authToken),
	}
}

func (p *PagerDutyClient) ListSchedules() ([]pagerduty.Schedule, error) {
	var schedules []pagerduty.Schedule
	var opts pagerduty.ListSchedulesOptions
	more := true
	for more {
		listSchedulesResponse, err := p.ApiClient.ListSchedules(opts)
		if err != nil {
			return nil, err
		}
		for _, schedule := range listSchedulesResponse.Schedules {
			schedules = append(schedules, schedule)
		}
		more = listSchedulesResponse.More
		opts.Offset = listSchedulesResponse.Limit
	}

	return schedules, nil
}

func (p *PagerDutyClient) ListServices(teamID string) ([]pagerduty.Service, error) {
	var opts pagerduty.ListServiceOptions
	opts.TeamIDs = []string{teamID}
	listServicesResponse, err := p.ApiClient.ListServices(opts)
	if err != nil {
		return nil, err
	}

	return listServicesResponse.Services, nil
}

func (p *PagerDutyClient) GetSchedule(scheduleID, startDate, endDate string) (*pagerduty.Schedule, error) {
	var opts pagerduty.GetScheduleOptions
	opts.Since = startDate
	opts.Until = endDate
	scheduleResponse, err := p.ApiClient.GetSchedule(scheduleID, opts)
	if err != nil {
		return nil, err
	}
	return scheduleResponse, nil
}
