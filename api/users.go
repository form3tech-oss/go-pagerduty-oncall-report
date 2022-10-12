package api

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
)

type User struct {
	ID       string
	Summary  string
	Name     string
	Email    string
	Timezone string
	Teams    []Team
}

func (p *PagerDutyClient) ListUsers() ([]*User, error) {
	var opts pagerduty.ListUsersOptions
	var userList []*User

	more := true
	for more {
		listUsersResponse, err := p.ApiClient.ListUsers(opts)
		if err != nil {
			return nil, err
		}

		for _, user := range listUsersResponse.Users {
			userList = append(userList, convertUser(&user))
		}
		more = listUsersResponse.More
		opts.Offset += listUsersResponse.Limit
	}

	return userList, nil
}

func (p *PagerDutyClient) GetUserById(id string) (*User, error) {
	pdUser, err := p.ApiClient.GetUser(id, pagerduty.GetUserOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by id (%s): %w", id, err)
	}

	return convertUser(pdUser), nil
}

func convertUser(user *pagerduty.User) *User {
	var userTeams []Team
	for _, team := range user.Teams {
		userTeams = append(userTeams, Team{
			ID:   team.ID,
			Name: team.Name,
		})
	}

	return &User{
		ID:       user.ID,
		Summary:  user.Summary,
		Name:     user.Name,
		Email:    user.Email,
		Timezone: user.Timezone,
		Teams:    userTeams,
	}
}
