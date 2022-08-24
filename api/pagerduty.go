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
