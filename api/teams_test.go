package api

import (
	"errors"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_ListTeams(t *testing.T) {
	tests := []struct {
		name        string
		clientSetup func(*clientMock)
		want        []Team
		wantErr     bool
	}{
		{
			name: "Failed to get list of teams",
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("ListTeams", mock.Anything).Once().Return(
					nil, errors.New("failed to get list of teams"))
			},
			wantErr: true,
		},
		{
			name: "Succesfully get list of teams",
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("ListTeams", mock.Anything).Once().Return(
					&pagerduty.ListTeamResponse{
						Teams: []pagerduty.Team{
							{
								APIObject: pagerduty.APIObject{
									ID: "QWERTY",
								},
								Name:        "Team 1",
								Description: "This is the team 1",
							},
						},
					}, nil)
			},
			want: []Team{
				{
					ID:          "QWERTY",
					Name:        "Team 1",
					Description: "This is the team 1",
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
			teamList, err := pdClient.ListTeams()
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			for i, wantTeam := range tt.want {
				assert.IsType(t, Team{}, teamList[i])
				assert.Equal(t, wantTeam.ID, teamList[i].ID)
				assert.Equal(t, wantTeam.Name, teamList[i].Name)
				assert.Equal(t, wantTeam.Description, teamList[i].Description)
			}
		})
	}
}
