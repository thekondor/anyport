package anyport

import (
	"net"
)

type FakeNetListener net.TCPAddr

func newFakeNetListener(port int) FakeNetListener {
	return FakeNetListener{Port: port}
}

func (fake FakeNetListener) Accept() (net.Conn, error) {
	panic("Should not be called")
}
func (fake FakeNetListener) Close() error {
	panic("Should not be called")
}
func (fake FakeNetListener) Addr() net.Addr {
	tcpAddr := net.TCPAddr(fake)
	return &tcpAddr
}
func (fake FakeNetListener) Network() string {
	panic("Should not be called")
}
func (fake FakeNetListener) String() string {
	panic("Should not be called")
}
