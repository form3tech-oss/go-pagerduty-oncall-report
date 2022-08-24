package cmd

import (
	"errors"
	"testing"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"

	"github.com/stretchr/testify/require"
)

func Test_listUsers(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*clientMock)
		wantErr   bool
	}{
		{
			name: "Successfully list pagerduty users",
			mockSetup: func(mock *clientMock) {
				mock.On("ListUsers").Return([]*api.User{
					{
						ID:       "QWERTY",
						Name:     "John Doe",
						Email:    "john.doe@email.com",
						Timezone: "Europe/London",
						Teams: []api.Team{
							{
								ID: "QWERTY",
							},
						},
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed to list pagerduty users",
			mockSetup: func(mock *clientMock) {
				mock.On("ListUsers").Return(nil, errors.New("failed to list"))
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
			err := pd.listUsers()
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
