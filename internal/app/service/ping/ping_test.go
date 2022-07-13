package ping

import (
	"context"
	"errors"
	mock_ping_service "github.com/ChristinaFomenko/shortener/internal/app/service/ping/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_service_Ping(t *testing.T) {
	tests := []struct {
		name string
		err  error
		exp  bool
	}{
		{
			name: "success",
			err:  nil,
			exp:  true,
		},
		{
			name: "repo err",
			err:  errors.New("test err"),
			exp:  false,
		},
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repoMock := mock_ping_service.NewMockurlRepo(ctrl)
		repoMock.EXPECT().Ping(ctx).Return(tt.err)

		s := NewService(repoMock)
		act := s.Ping(ctx)

		assert.Equal(t, tt.exp, act)
	}
}
