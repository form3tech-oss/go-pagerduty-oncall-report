package api

import (
	"errors"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_GetUserById(t *testing.T) {
	tests := []struct {
		name        string
		clientSetup func(*clientMock)
		want        User
		wantErr     bool
	}{
		{
			name: "Successfully get user by ID",
			want: User{
				ID:       "QWERTY",
				Name:     "John Doe",
				Email:    "john.doe@email.com",
				Timezone: "Europe/London",
				Teams: []Team{
					{
						ID: "QWERTY",
					},
				},
			},
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("GetUser", mock.Anything, mock.Anything).Once().Return(&pagerduty.User{
					APIObject: pagerduty.APIObject{
						ID: "QWERTY",
					},
					Name:     "John Doe",
					Email:    "john.doe@email.com",
					Timezone: "Europe/London",
					Teams: []pagerduty.Team{
						{
							APIObject: pagerduty.APIObject{
								ID: "QWERTY",
							},
						},
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed get user by ID",
			want: User{
				ID:       "QWERTY",
				Name:     "John Doe",
				Email:    "john.doe@email.com",
				Timezone: "Europe/London",
			},
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("GetUser", mock.Anything, mock.Anything).Once().Return(
					nil, errors.New("failed to get user by id"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockedClient := &clientMock{}
			if tt.clientSetup != nil {
				tt.clientSetup(mockedClient)
			}

			pdClient := PagerDutyClient{ApiClient: mockedClient}
			user, err := pdClient.GetUserById("randomID")
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, user.ID, tt.want.ID)
			assert.Equal(t, user.Name, tt.want.Name)
			assert.Equal(t, user.Email, tt.want.Email)
			assert.Equal(t, user.Timezone, tt.want.Timezone)
			assert.Equal(t, user.Teams[0].ID, tt.want.Teams[0].ID)
		})
	}
}

func Test_ListUsers(t *testing.T) {
	tests := []struct {
		name        string
		clientSetup func(*clientMock)
		want        []*User
		wantErr     bool
	}{
		{
			name: "Failed to get list of users",
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("ListUsers", mock.Anything).Once().Return(
					nil, errors.New("failed to get list of users"))
			},
			wantErr: true,
		},
		{
			name: "Successfully get list of users",
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("ListUsers", mock.Anything).Once().Return(
					&pagerduty.ListUsersResponse{
						Users: []pagerduty.User{
							{
								APIObject: pagerduty.APIObject{
									ID: "QWERTY",
								},
								Name:     "John Doe",
								Email:    "john.doe@email.com",
								Timezone: "Europe/London",
								Teams: []pagerduty.Team{
									{
										APIObject: pagerduty.APIObject{
											ID: "QWERTY",
										},
									},
								},
							},
							{
								APIObject: pagerduty.APIObject{
									ID: "QWERTY2",
								},
								Name:     "Jane Doe",
								Email:    "jane.doe@email.com",
								Timezone: "Europe/London",
								Teams: []pagerduty.Team{
									{
										APIObject: pagerduty.APIObject{
											ID: "QWERTY2",
										},
									},
								},
							},
						},
					}, nil)
			},
			want: []*User{
				{
					ID:       "QWERTY",
					Name:     "John Doe",
					Email:    "john.doe@email.com",
					Timezone: "Europe/London",
					Teams: []Team{
						{
							ID: "QWERTY",
						},
					},
				},
				{
					ID:       "QWERTY2",
					Name:     "Jane Doe",
					Email:    "jane.doe@email.com",
					Timezone: "Europe/London",
					Teams: []Team{
						{
							ID: "QWERTY2",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedClient := &clientMock{}
			if tt.clientSetup != nil {
				tt.clientSetup(mockedClient)
			}

			pdClient := PagerDutyClient{ApiClient: mockedClient}
			userList, err := pdClient.ListUsers()
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			for i, wantUser := range tt.want {
				assert.IsType(t, User{}, userList[i])
				assert.Equal(t, wantUser.ID, userList[i].ID)
				assert.Equal(t, wantUser.Name, userList[i].Name)
				assert.Equal(t, wantUser.Email, userList[i].Email)
				assert.Equal(t, wantUser.Timezone, userList[i].Timezone)
				assert.Equal(t, wantUser.Teams[0].ID, userList[i].Teams[0].ID)
			}

			require.NoError(t, err)

		})
	}
}
