package cmd

import (
	"errors"
	"testing"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_pagerDutyClient_loadUsersInMemoryCache(t *testing.T) {
	users := []*api.User{
		{
			ID:   "1",
			Name: "John Doe",
		},
		{
			ID:   "2",
			Name: "Mary Jane",
		},
	}

	tests := []struct {
		name      string
		mockSetup func(*clientMock)
		want      []*api.User
		wantErr   bool
	}{
		{
			name: "Successfully load users in memory",
			mockSetup: func(mock *clientMock) {
				mock.On("ListUsers").Once().Return(users, nil)
			},
			want:    users,
			wantErr: false,
		},
		{
			name: "Failed load users in memory",
			mockSetup: func(mock *clientMock) {
				mock.On("ListUsers").Once().Return(nil, errors.New("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedClient := &clientMock{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockedClient)
			}

			pd := pagerDutyClient{client: mockedClient}
			err := pd.loadUsersInMemoryCache()
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, users, pd.cachedUsers)
		})
	}
}

func Test_pagerDutyClient_getUserTimezone(t *testing.T) {
	tests := []struct {
		name        string
		cachedUsers []*api.User
		mockSetup   func(mock *clientMock)
		want        string
		wantErr     bool
	}{
		{
			name: "Successfully find the user timezone",
			cachedUsers: []*api.User{
				{
					ID:       "USER_ID",
					Timezone: "Europe/London",
				},
			},
			want:    "Europe/London",
			wantErr: false,
		},
		{
			name: "User with empty timezone will fail to return timezone",
			cachedUsers: []*api.User{
				{
					ID: "USER_ID",
				},
			},
			want:    "Europe/London",
			wantErr: true,
		},
		{
			name: "User not cached will fail to return timezone",
			cachedUsers: []*api.User{
				{
					ID: "NOT_THE_USER_ID",
				},
			},
			want:    "Europe/London",
			wantErr: true,
		},
		{
			name: "If user not cached it will load users in cache and successfully return timezone",
			mockSetup: func(mock *clientMock) {
				mock.On("ListUsers").Once().Return([]*api.User{
					{ID: "USER_ID", Timezone: "Europe/London"},
				}, nil)
			},
			want:    "Europe/London",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedClient := &clientMock{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockedClient)
			}

			pd := pagerDutyClient{client: mockedClient, cachedUsers: tt.cachedUsers}

			got, err := pd.getUserTimezone("USER_ID")
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
