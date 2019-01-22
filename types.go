package anyport

import (
	"crypto/tls"
	"net"
)

// AnyPort structure is used to keep information about active port being listened.
type AnyPort struct {
	// Active listener bound to an available port
	Listener net.Listener
	// Port number. Shortcut for 'Listener.Addr().(*net.TCPAddr)' expression.
	PortNumber int
}

// Listener interface is used to abstract net.Listen() and tls.Listen() calls for the sake of testing; as well as for custom Listen() routine implementations.
type Listener interface {
	ListenInsecure(network, address string) (net.Listener, error)
	ListenSecure(network, address string, config *tls.Config) (net.Listener, error)
}
