package cmd

import (
	"errors"
	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"
	"testing"
	"time"

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
		name                string
		cachedUsers         []*api.User
		defaultUserTimezone string
		mockSetup           func(mock *clientMock)
		want                string
		wantErr             bool
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
			name: "User with empty timezone will receive default configured timezone",
			cachedUsers: []*api.User{
				{
					ID: "USER_ID",
				},
			},
			defaultUserTimezone: "Europe/London",
			want:                "Europe/London",
			wantErr:             false,
		},
		{
			name: "User not cached will receive default configured timezone",
			cachedUsers: []*api.User{
				{
					ID: "NOT_THE_USER_ID",
				},
			},
			defaultUserTimezone: "America/Sao_Paulo",
			want:                "America/Sao_Paulo",
			wantErr:             false,
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

			pd := pagerDutyClient{
				client:              mockedClient,
				cachedUsers:         tt.cachedUsers,
				defaultUserTimezone: tt.defaultUserTimezone,
			}

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

func Test_pagerDutyClient_convertToUserLocalTimezone(t *testing.T) {
	tests := []struct {
		name         string
		scheduleDate string
		cachedUsers  []*api.User
		mockSetup    func(mock *clientMock)
		want         string
		wantErr      bool
	}{
		{
			name:         "Successfully converts schedule BST date to user's EDT local date",
			scheduleDate: "01 Sep 22 17:00 BST",
			cachedUsers: []*api.User{
				{
					ID:       "USER_ID",
					Timezone: "America/New_York",
				},
			},
			want:    "01 Sep 22 12:00 EDT",
			wantErr: false,
		},
		{
			name:         "Fails to list users fails to convert user timezone",
			scheduleDate: "01 Sep 22 17:00 BST",
			mockSetup: func(mock *clientMock) {
				mock.On("ListUsers").Once().Return(nil, errors.New("failed"))
			},
			wantErr: true,
		},
		{
			name:         "Fails to convert due non existent timezone",
			scheduleDate: "01 Sep 22 17:00 BST",
			cachedUsers: []*api.User{
				{
					ID:       "USER_ID",
					Timezone: "Mars/Kaiser_Sea",
				},
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

			pd := pagerDutyClient{
				client:      mockedClient,
				cachedUsers: tt.cachedUsers,
			}

			scheduleDate, err := time.Parse(time.RFC822, tt.scheduleDate)
			require.NoError(t, err)

			got, err := pd.convertToUserLocalTimezone(scheduleDate, "USER_ID")
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			expectedDate, err := time.Parse(time.RFC822, tt.want)
			require.NoError(t, err)

			assert.Equal(t, expectedDate.Hour(), got.Hour())
			assert.Equal(t, expectedDate.Minute(), got.Minute())
		})
	}
}
