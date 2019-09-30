package servicemanager

import (
	"encoding/json"
	"gotest.tools/assert"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestHomeHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	homeHandler().ServeHTTP(w, r)

	if !strings.Contains(w.Body.String(), ServiceName) {
		t.Fatal("home should contain service name")
	}

	if w.Code != 200 {
		t.Fatal("home should return 200 code")
	}
}

type mockService struct {
	service ServiceInfo
}

func (m *mockService) RegisterService(service ServiceInfo) {
	m.service = service
}

func TestRegisterHandler(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/?service=Api%20endpoint&endpoint=http://api:9000&value=123", nil)
		m := mockService{}
		registerHandler(&m).ServeHTTP(w, r)

		assert.Equal(t, m.service, ServiceInfo{
			Name:            "Api endpoint",
			Endpoint:        "http://api:9000",
			HealthCheckTime: time.Now().Unix(),
			Value:           "123",
		}, "The service registered dont match")
	})

	t.Run("POST", func(t *testing.T) {
		w := httptest.NewRecorder()

		data := url.Values{}
		data.Add("value", "my-value-from-post")

		r := httptest.NewRequest("POST", "/?service=Api%20endpoint&endpoint=http://api:9000", strings.NewReader(data.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		m := mockService{}
		registerHandler(&m).ServeHTTP(w, r)

		assert.Equal(t, m.service, ServiceInfo{
			Name:            "Api endpoint",
			Endpoint:        "http://api:9000",
			HealthCheckTime: time.Now().Unix(),
			Value:           "my-value-from-post",
		}, "The service registered dont match")
	})
}

func TestServicesHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	services := []ServiceInfo{
		{
			Name:            "TestService1",
			Endpoint:        "http://test-service:9000",
			HealthCheckTime: time.Now().Unix(),
			Value:           "123",
		},
		{
			Name:            "TestService2",
			Endpoint:        "http://test-service2:9000",
			HealthCheckTime: time.Now().Unix(),
			Value:           "321",
		},
	}

	servicesFun := func() []ServiceInfo {
		return services
	}

	servicesHandler(servicesFun).ServeHTTP(w, r)

	data, err := json.Marshal(&services)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(data) + "\n", w.Body.String(), "Services dont match")
}
