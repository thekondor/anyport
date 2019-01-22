package anyport

import (
	"fmt"
	"net"
	"strconv"
)

type listenFunc func(host string) (net.Listener, error)

func dispatchListen(addr string, listen listenFunc) (AnyPort, error) {
	listener, err := listen(addr)
	if nil != err {
		return AnyPort{}, err
	}

	return AnyPort{Listener: listener, PortNumber: listener.Addr().(*net.TCPAddr).Port}, nil
}

func dispatchRangeListen(host string, minPort, maxPort int, listen listenFunc) (AnyPort, error) {
	for port := minPort; port <= maxPort; port++ {
		addr := fmt.Sprintf("%s:%d", host, port)
		anyPort, err := dispatchListen(addr, listen)
		if nil == err {
			return anyPort, nil
		}
	}

	return AnyPort{}, fmt.Errorf("Failed to bind '%s' within [%d:%d]", host, minPort, maxPort)
}

func parseRange(rawRange []string) (min int, max int, err error) {
	if 2 != len(rawRange) {
		panic("Internal error. Invalid range array size")
	}

	min, err = strconv.Atoi(rawRange[0])
	if nil != err {
		return -1, -1, fmt.Errorf("Invalid min port: %s", rawRange[0])
	}

	max, err = strconv.Atoi(rawRange[1])
	if nil != err {
		return -1, -1, fmt.Errorf("Invalid max port: %s", rawRange[1])
	}

	return
}
