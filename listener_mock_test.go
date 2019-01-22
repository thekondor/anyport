package anyport

import (
	"crypto/tls"
	"github.com/stretchr/testify/mock"
	"net"
)

type ListenerMock struct {
	mock.Mock
}

func (mock *ListenerMock) ListenInsecure(network, address string) (net.Listener, error) {
	args := mock.Called(network, address)
	return args.Get(0).(net.Listener), args.Error(1)
}

func (mock *ListenerMock) ListenSecure(network, address string, config *tls.Config) (net.Listener, error) {
	args := mock.Called(network, address, config)
	return args.Get(0).(net.Listener), args.Error(1)
}
