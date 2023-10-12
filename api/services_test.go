package api

import (
	"errors"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_ListServices(t *testing.T) {
	tests := []struct {
		name        string
		clientSetup func(*clientMock)
		want        []*Service
		wantErr     bool
	}{
		{
			name: "Failed to get list of services",
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("ListServices", mock.Anything).Once().Return(
					nil, errors.New("failed to get list of services"))
			},
			wantErr: true,
		},
		{
			name: "Successfully get list of services",
			clientSetup: func(clientMock *clientMock) {
				clientMock.On("ListServices", mock.Anything).Once().Return(
					&pagerduty.ListServiceResponse{
						Services: []pagerduty.Service{
							{
								APIObject: pagerduty.APIObject{
									ID: "QWERTY",
								},
								Name:        "Service 1",
								Description: "This is the service 1",
							},
						},
					}, nil)
			},
			want: []*Service{
				{
					ID:   "QWERTY",
					Name: "Service 1",
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
			serviceList, err := pdClient.ListServices("QWERTY")
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			for i, wantService := range tt.want {
				assert.IsType(t, &Service{}, serviceList[i])
				assert.Equal(t, wantService.ID, serviceList[i].ID)
				assert.Equal(t, wantService.Name, serviceList[i].Name)
			}
		})
	}
}
