package anyport

import (
	"crypto/tls"
	"net"
)

type stdListener struct{}

func (listener stdListener) ListenInsecure(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

func (listener stdListener) ListenSecure(network, address string, config *tls.Config) (net.Listener, error) {
	return tls.Listen(network, address, config)
}
