package cmd

import (
	"errors"
	"testing"

	"github.com/form3tech-oss/go-pagerduty-oncall-report/api"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_listServices(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*clientMock)
		wantErr   bool
	}{
		{
			name: "Successfully list pagerduty services",
			mockSetup: func(clientMock *clientMock) {
				clientMock.On("ListServices", mock.Anything).Return([]*api.Service{
					{
						ID:   "QWERTY",
						Name: "Service 1",
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed to list pagerduty services",
			mockSetup: func(clientMock *clientMock) {
				clientMock.On("ListServices", mock.Anything).Return(nil, errors.New("failed to list"))
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
			err := pd.listServices("fake-service-id")
			mockedClient.AssertExpectations(t)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
