package anyport

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net"
	"testing"
)

type AnyPortTest struct {
	suite.Suite
	originalListener Listener
	listenerMock     ListenerMock
}

var (
	invalidNetListener net.Listener = newFakeNetListener(-100)
	fakeNetListener    net.Listener = newFakeNetListener(102)
)
var (
	dummyTlsConfig tls.Config
)

func TestAnyPort(t *testing.T) {
	suite.Run(t, new(AnyPortTest))
}

func (test *AnyPortTest) SetupTest() {
	test.originalListener = CurrentListenerInstance
	test.listenerMock = ListenerMock{}
	CurrentListenerInstance = &test.listenerMock
}

func (test *AnyPortTest) TearDownTest() {
	CurrentListenerInstance = test.originalListener
}

func (test *AnyPortTest) Test_ListenInsecure_Errors_OnUnderlyingError() {
	test.listenerMock.On("ListenInsecure", "tcp", "999.0.0.0:0").
		Return(invalidNetListener, errors.New("listen error"))

	_, err := ListenInsecure("999.0.0.0")
	test.Assert().Error(err)
}

func (test *AnyPortTest) Test_ListenSecure_Errors_OnUnderlyingError() {
	test.listenerMock.On("ListenSecure", "tcp", "999.0.0.0:0", &dummyTlsConfig).
		Return(invalidNetListener, errors.New("listen error"))

	_, err := ListenSecure("999.0.0.0", &dummyTlsConfig)
	test.Assert().Error(err)
}

func (test *AnyPortTest) Test_ListenInsecure_Succeeds_OnUnderlyingSuccess() {
	test.listenerMock.On("ListenInsecure", "tcp", "0.0.0.0:0").
		Return(fakeNetListener, nil)

	anyPort, err := ListenInsecure("0.0.0.0")
	test.Require().NoError(err)
	test.Assert().Equal(anyPort.Listener, fakeNetListener)
	test.Assert().Equal(102, anyPort.PortNumber)
}

func (test *AnyPortTest) Test_ListenSecure_Succeeds_OnUnderlyingSuccess() {
	test.listenerMock.On("ListenSecure", "tcp", "0.0.0.0:0", &dummyTlsConfig).
		Return(fakeNetListener, nil)

	anyPort, err := ListenSecure("0.0.0.0", &dummyTlsConfig)
	test.Require().NoError(err)
	test.Assert().Equal(anyPort.Listener, fakeNetListener)
	test.Assert().Equal(102, anyPort.PortNumber)
}

func (test *AnyPortTest) Test_ListenRange_Errors_OnNoAvailablePortFound() {
	anyAddr := func(string) bool { return true }

	test.T().Run("Insecure", func(t *testing.T) {
		test.listenerMock.On("ListenInsecure", "tcp", mock.MatchedBy(anyAddr)).
			Return(invalidNetListener, errors.New("Insecure port is not available"))
		_, err := ListenInsecure("0.0.0.0:100-200")
		assert.Error(t, err)
	})

	test.T().Run("Secure", func(t *testing.T) {
		test.listenerMock.On("ListenSecure", "tcp", mock.MatchedBy(anyAddr), &dummyTlsConfig).
			Return(invalidNetListener, errors.New("Secure port is not available"))

		_, err := ListenSecure("0.0.0.0:200-300", &dummyTlsConfig)
		assert.Error(t, err)
	})
}

func (test *AnyPortTest) Test_ListenRange_Errors_OnMalformedPortsRange() {
	cases := []string{"123-", "a-123", "123-b", "999-1"}

	sutFunctions := map[string]func(addr string) (AnyPort, error){
		"Insecure": func(addr string) (AnyPort, error) {
			return ListenInsecure(addr)
		},
		"Secure": func(addr string) (AnyPort, error) {
			return ListenSecure(addr, &tls.Config{})
		},
	}

	for sutName, sutFn := range sutFunctions {
		for _, input := range cases {
			testCaseName := fmt.Sprintf("%s(), case:%s", sutName, input)
			test.T().Run(testCaseName, func(t *testing.T) {
				_, err := sutFn("0.0.0.0:" + input)
				assert.Error(t, err)
			})
		}
	}
}

func (test *AnyPortTest) Test_ListenInsecureRange_Succeeds_OnUnderlyingSuccess() {
	test.listenerMock.
		On("ListenInsecure", "tcp", "0.0.0.0:100").Return(invalidNetListener, errors.New("port already in use")).Once().
		On("ListenInsecure", "tcp", "0.0.0.0:101").Return(invalidNetListener, errors.New("permission denied")).Once().
		On("ListenInsecure", "tcp", "0.0.0.0:102").Return(fakeNetListener, nil).Once()

	anyPort, err := ListenInsecure("0.0.0.0:100-200")
	test.Require().NoError(err)
	test.Assert().Equal(102, anyPort.PortNumber)
}

func (test *AnyPortTest) Test_ListenSecureRange_Succeeds_OnUnderlyingSuccess() {
	test.listenerMock.
		On("ListenSecure", "tcp", "0.0.0.0:100", &dummyTlsConfig).Return(invalidNetListener, errors.New("port already in use")).Once().
		On("ListenSecure", "tcp", "0.0.0.0:101", &dummyTlsConfig).Return(invalidNetListener, errors.New("permission denied")).Once().
		On("ListenSecure", "tcp", "0.0.0.0:102", &dummyTlsConfig).Return(fakeNetListener, nil).Once()

	anyPort, err := ListenSecure("0.0.0.0:100-200", &dummyTlsConfig)
	test.Require().NoError(err)
	test.Assert().Equal(102, anyPort.PortNumber)
}
