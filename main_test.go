// cmd/myapp/main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theCompanyDream/pingtest/internal"
)

// MockCommands is a mock implementation of the Commands interface
type MockCommands struct {
	mock.Mock
}

func (m *MockCommands) Ping(host string) (commands.PingResult, error) {
	args := m.Called(host)
	return args.Get(0).(commands.PingResult), args.Error(1)
}

func (m *MockCommands) GetSystemInfo() (commands.SystemInfo, error) {
	args := m.Called()
	return args.Get(0).(commands.SystemInfo), args.Error(1)
}

// TestHandlePing_Success tests the /ping endpoint with a successful ping
func TestHandlePing_Success(t *testing.T) {
	// Create a mock Commands
	mockCmds := new(MockCommands)

	// Define expected host and result
	host := "google.com"
	expectedResult := commands.PingResult{
		Successful: true,
		Time:       123 * time.Millisecond,
		Packets:    4,
	}

	// Setup expectations
	mockCmds.On("Ping", host).Return(expectedResult, nil)

	// Create handler with mock Commands
	handler := makeHandlePing(mockCmds)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/ping?host="+host, nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response commands.PingResult
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected valid JSON response")

	assert.Equal(t, expectedResult.Successful, response.Successful, "Successful field mismatch")
	assert.Equal(t, expectedResult.Time, response.Time, "Time field mismatch")
	assert.Equal(t, expectedResult.Packets, response.Packets, "Packets field mismatch")

	// Assert that the expectations were met
	mockCmds.AssertExpectations(t)
}

// TestHandlePing_MissingHost tests the /ping endpoint without the host parameter
func TestHandlePing_MissingHost(t *testing.T) {
	// Create a mock Commands
	mockCmds := new(MockCommands)

	// Create handler with mock Commands
	handler := makeHandlePing(mockCmds)

	// Create a test HTTP request without the host parameter
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status Bad Request")
	assert.Equal(t, "Host parameter not found in query string\n", rr.Body.String(), "Unexpected response body")

	// No calls should have been made to Ping
	mockCmds.AssertNotCalled(t, "Ping", mock.Anything)
}

// TestHandlePing_Failure tests the /ping endpoint with a failed ping
func TestHandlePing_Failure(t *testing.T) {
	// Create a mock Commands
	mockCmds := new(MockCommands)

	// Define expected host and error
	host := "invalid_host"
	expectedError := errors.New("ping failed")

	// Setup expectations
	mockCmds.On("Ping", host).Return(commands.PingResult{
		Successful: false,
		Time:       0,
		Packets:    0,
	}, expectedError)

	// Create handler with mock Commands
	handler := makeHandlePing(mockCmds)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/ping?host="+host, nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status Bad Request")
	assert.Equal(t, "Error Running Ping Request\n", rr.Body.String(), "Unexpected response body")

	// Assert that Ping was called with the correct host
	mockCmds.AssertCalled(t, "Ping", host)
}

// TestHandleGetSystemInfo_Success tests the / endpoint with successful system info retrieval
func TestHandleGetSystemInfo_Success(t *testing.T) {
	// Create a mock Commands
	mockCmds := new(MockCommands)

	// Define expected system info
	expectedInfo := commands.SystemInfo{
		Hostname:  "test-host",
		IPAddress: "192.168.1.100",
	}

	// Setup expectations
	mockCmds.On("GetSystemInfo").Return(expectedInfo, nil)

	// Create handler with mock Commands
	handler := makeHandleGetSystemInfo(mockCmds)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response commands.SystemInfo
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected valid JSON response")

	assert.Equal(t, expectedInfo.Hostname, response.Hostname, "Hostname mismatch")
	assert.Equal(t, expectedInfo.IPAddress, response.IPAddress, "IPAddress mismatch")

	// Assert that the expectations were met
	mockCmds.AssertExpectations(t)
}

// TestHandleGetSystemInfo_Failure tests the / endpoint with a failure in system info retrieval
func TestHandleGetSystemInfo_Failure(t *testing.T) {
	// Create a mock Commands
	mockCmds := new(MockCommands)

	// Define expected error
	expectedError := errors.New("failed to retrieve system info")

	// Setup expectations
	mockCmds.On("GetSystemInfo").Return(commands.SystemInfo{}, expectedError)

	// Create handler with mock Commands
	handler := makeHandleGetSystemInfo(mockCmds)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
	assert.Equal(t, "Error retrieving system info\n", rr.Body.String(), "Unexpected response body")

	// Assert that GetSystemInfo was called
	mockCmds.AssertCalled(t, "GetSystemInfo")
}
