// Package anyport provides with a facility to bind any available or random port to listen for TCP connection.
//
// 'Listen*()' functions bind TCP port for the address specified:
//   myhost.mydomain:1234
// will bind port 1234 if available.
// When no port is specified within address argument, a random port is bound if available. E.g.
//   myhost.mydomain
// When port is specified in 'A-B' a way, first available port in [A, B] range is bound. E.g.:
//   myhost.mydomain:123-456
// will bind for listening any available port in range [123-456].
// When a requested port is not available or could not be found, a corresponding error is returned.
//
// No extra logic over standard 'Listen()' is introduced.
package anyport

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
)

// Default implementation of Listener interface built over standard 'net.Listen()' and 'tls.Listen()' calls. More likely should never be changed.
var CurrentListenerInstance Listener = stdListener{}

func listen(addr string, listen func(host string) (net.Listener, error)) (AnyPort, error) {
	addrComponents := strings.Split(addr, ":")
	if len(addrComponents) < 2 {
		return dispatchListen(fmt.Sprintf("%s:0", addr), listen)
	}

	portsRange := strings.Split(addrComponents[1], "-")
	if len(portsRange) < 2 {
		return dispatchListen(addr, listen)
	}

	minPort, maxPort, err := parseRange(portsRange)
	if nil != err {
		return AnyPort{}, fmt.Errorf("Invalid min-max port range, could not be parsed: %+v", portsRange)
	}

	if maxPort < minPort {
		return AnyPort{}, fmt.Errorf("Invalid min-max port range, min and max are incorrect: %d:%d", minPort, maxPort)
	}

	return dispatchRangeListen(addrComponents[0], minPort, maxPort, listen)
}

// Binds TCP port for address speciefied. Plain connection.
func ListenInsecure(addr string) (AnyPort, error) {
	return listen(addr, func(adjustedAddr string) (net.Listener, error) {
		return CurrentListenerInstance.ListenInsecure("tcp", adjustedAddr)
	})
}

// Binds TCP port for address specified. Connection over TLS.
func ListenSecure(addr string, config *tls.Config) (AnyPort, error) {
	return listen(addr, func(adjustedAddr string) (net.Listener, error) {
		return CurrentListenerInstance.ListenSecure("tcp", adjustedAddr, config)
	})
}
