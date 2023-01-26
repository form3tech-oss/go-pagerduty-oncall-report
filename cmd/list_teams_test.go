package cmd

import (
	"errors"
	"testing"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"

	"github.com/stretchr/testify/require"
)

func Test_listTeams(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*clientMock)
		wantErr   bool
	}{
		{
			name: "Successfully list pagerduty teams",
			mockSetup: func(mock *clientMock) {
				mock.On("ListTeams").Return([]*api.Team{
					{
						ID:   "QWERTY",
						Name: "Team 1",
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed to list pagerduty teams",
			mockSetup: func(mock *clientMock) {
				mock.On("ListTeams").Return(nil, errors.New("failed to list"))
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
			err := pd.listTeams()
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
