package cmd

import (
	"errors"
	"testing"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_listSchedules(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*clientMock)
		wantErr   bool
	}{
		{
			name: "Successfully list pagerduty schedules",
			mockSetup: func(clientMock *clientMock) {
				clientMock.On("ListSchedules", mock.Anything).Once().Return([]*api.Schedule{
					{
						ID:          "QWERTY",
						Name:        "Schedule 1",
						TimeZone:    "Europe/London",
						Description: "This is the schedule 1",
						FinalSchedule: api.ScheduleLayer{
							ID: "QWERTY1",
							RenderedScheduleEntries: []api.RenderedScheduleEntry{
								{
									Start: "2022-08-24T09:35:12+01:00",
								},
							},
						},
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed to list pagerduty schedules",
			mockSetup: func(clientMock *clientMock) {
				clientMock.On("ListSchedules", mock.Anything).Once().Return(nil, errors.New("failed to list"))
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
			err := pd.listSchedules()
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
