package servicemanager

import (
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	var err error
	_, err = New(time.Second, []string{})
	assert.NilError(t, err, "should not return error when passing empty webhook urls")

	_, err = New(time.Second, []string{
		"http://valid.com:9000",
	})
	assert.NilError(t, err, "should not return error when passing valid webhook url")

	_, err = New(time.Second, []string{
		"http://192.168.0.10:9000",
	})
	assert.NilError(t, err, "should not return error when passing valid ip webhook url")

	_, err = New(time.Second, []string{
		"http://valid:9000",
	})
	assert.NilError(t, err, "should not return error when passing valid docker webhook url")

	_, err = New(time.Second, []string{
		"inva lid",
	})
	assert.ErrorContains(t, err, "invalid webhook url", "should return error when passing a invalid url (single name)")

	_, err = New(time.Second, []string{
		"invalid",
	})
	assert.ErrorContains(t, err, "invalid webhook url", "should return error when passing a invalid url (spaces)")
}

func TestServer_GetServices(t *testing.T) {
	s, _ := New(time.Second, []string{})
	assert.Equal(t, len(s.GetServices()), 0, "Returns 0 services")
}

func TestServer_RegisterService(t *testing.T) {
	s, _ := New(time.Second, []string{})
	assert.Equal(t, len(s.GetServices()), 0, "Returns 0 services")

	s.RegisterService(ServiceInfo{
		Name:            "test-service",
		Endpoint:        "http://test-service:9000",
		HealthCheckTime: time.Now().Unix(),
		Value:           nil,
	})

	assert.Equal(t, len(s.GetServices()), 1, "Returns 1 services")

	time.Sleep(time.Second)
	assert.Equal(t, len(s.GetServices()), 0, "Returns 0 services")
}

func TestServer_checkForServiceChanges(t *testing.T) {

}