package api

import (
	"errors"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_ListSchedules(t *testing.T) {
	tests := []struct {
		name        string
		clientSetup func(*clientMock)
		want        []*Schedule
		wantErr     bool
	}{
		{
			name: "Failed to get list of schedules",
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("ListSchedules", mock.Anything).Once().Return(
					nil, errors.New("failed to get list of schedules"))
			},
			wantErr: true,
		},
		{
			name: "Successfully get list of schedules",
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("ListSchedules", mock.Anything).Once().Return(
					&pagerduty.ListSchedulesResponse{
						APIListObject: pagerduty.APIListObject{},
						Schedules: []pagerduty.Schedule{
							{
								APIObject: pagerduty.APIObject{
									ID: "QWERTY",
								},
								Name:        "Schedule 1",
								TimeZone:    "Europe/London",
								Description: "This is the schedule 1",
								FinalSchedule: pagerduty.ScheduleLayer{
									APIObject: pagerduty.APIObject{
										ID: "QWERTY1",
									},
									Name:  "Final Schedule 1",
									Start: "2022-08-24T09:35:12+01:00",
									RenderedScheduleEntries: []pagerduty.RenderedScheduleEntry{
										{
											Start: "2022-08-24T09:35:12+01:00",
										},
									},
								},
							},
						},
					}, nil)
			},
			want: []*Schedule{
				{
					ID:       "QWERTY",
					Name:     "Schedule 1",
					TimeZone: "Europe/London",
					FinalSchedule: ScheduleLayer{
						RenderedScheduleEntries: []RenderedScheduleEntry{
							{
								Start: "2022-08-24T09:35:12+01:00",
							},
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
			scheduleList, err := pdClient.ListSchedules()
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			for i, wantSchedule := range tt.want {
				assert.IsType(t, &Schedule{}, scheduleList[i])
				assert.Equal(t, wantSchedule.ID, scheduleList[i].ID)
				assert.Equal(t, wantSchedule.Name, scheduleList[i].Name)
				assert.Equal(t, wantSchedule.TimeZone, scheduleList[i].TimeZone)

				assert.IsType(t, ScheduleLayer{}, scheduleList[i].FinalSchedule)
				assert.IsType(t, []RenderedScheduleEntry{}, scheduleList[i].FinalSchedule.RenderedScheduleEntries)
			}
		})
	}
}

func Test_GetSchedule(t *testing.T) {
	tests := []struct {
		name        string
		clientSetup func(*clientMock)
		want        Schedule
		wantErr     bool
	}{
		{
			name: "Successfully get schedule by ID",
			want: Schedule{
				ID:       "QWERTY",
				Name:     "Schedule 1",
				TimeZone: "Europe/London",
				FinalSchedule: ScheduleLayer{
					RenderedScheduleEntries: []RenderedScheduleEntry{
						{
							Start: "2022-08-24T09:35:12+01:00",
						},
					},
				},
			},
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("GetSchedule", mock.Anything, mock.Anything).Once().Return(
					&pagerduty.Schedule{
						APIObject: pagerduty.APIObject{
							ID: "QWERTY",
						},
						Name:        "Schedule 1",
						TimeZone:    "Europe/London",
						Description: "This is the schedule 1",
						FinalSchedule: pagerduty.ScheduleLayer{
							APIObject: pagerduty.APIObject{
								ID: "QWERTY1",
							},
							Name:  "Final Schedule 1",
							Start: "2022-08-24T09:35:12+01:00",
							RenderedScheduleEntries: []pagerduty.RenderedScheduleEntry{
								{
									Start: "2022-08-24T09:35:12+01:00",
								},
							},
						},
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed get schedule by ID",
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("GetSchedule", mock.Anything, mock.Anything).Once().Return(
					nil, errors.New("failed to get schedule by id"))
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
			schedule, err := pdClient.GetSchedule("randomID", "randomStartDate", "randomEndDate")
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.IsType(t, &Schedule{}, schedule)
			assert.Equal(t, tt.want.ID, schedule.ID)
			assert.Equal(t, tt.want.Name, schedule.Name)

			assert.IsType(t, ScheduleLayer{}, schedule.FinalSchedule)
			assert.IsType(t, []RenderedScheduleEntry{}, schedule.FinalSchedule.RenderedScheduleEntries)
		})
	}
}
