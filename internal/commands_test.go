// commands_test.go
package internal

import (
	"os"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	probing "github.com/prometheus-community/pro-bing"
)

// MockPinger is a mock implementation of the probing.Pinger interface
type MockPinger struct {
	mock.Mock
}

func (m *MockPinger) SetPrivileged(privileged bool) {
	m.Called(privileged)
}

func (m *MockPinger) Run() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPinger) Statistics() *probing.Statistics {
	args := m.Called()
	return args.Get(0).(*probing.Statistics)
}

// TestPingSuccess tests the Ping function with a valid host
func TestPingSuccess(t *testing.T) {
	// Skip the test if not privileged (requires root)
	if os.Geteuid() != 0 {
		t.Skip("Skipping test as it requires root privileges")
	}

	host := "google.com"
	result, err := Ping(host)

	require.NoError(t, err, "Expected no error for a valid host")
	assert.True(t, result.Successful, "Ping should be successful")
	assert.Greater(t, result.Packets, 0, "At least one packet should be received")
	assert.LessOrEqual(t, result.Time, 10*time.Second, "Ping should complete within the timeout")
}

// TestPingInvalidHost tests the Ping function with an invalid host
func TestPingInvalidHost(t *testing.T) {
	host := "invalid_host"

	result, err := Ping(host)

	require.Error(t, err, "Expected an error for an invalid host")
	assert.False(t, result.Successful, "Ping should not be successful")
	assert.Equal(t, 0, result.Packets, "No packets should be received")
}

// TestGetSystemInfo tests the GetSystemInfo function
func TestGetSystemInfo(t *testing.T) {
	info, err := GetSystemInfo()

	require.NoError(t, err, "Expected no error when getting system info")
	assert.NotEmpty(t, info.Hostname, "Hostname should not be empty")
	assert.NotEmpty(t, info.IPAddress, "IP Address should not be empty")

	parsedIP := net.ParseIP(info.IPAddress)
	assert.NotNil(t, parsedIP, "IP Address should be a valid IP")
}
