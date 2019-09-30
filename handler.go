package servicemanager

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type httpResponseError struct {
	Service string `json:"service"`
	Error   string `json:"error"`
}

type httpResponseOK struct {
	Service string `json:"service"`
	OK      bool   `json:"ok"`
}

func newHTTPResponseError(err error) httpResponseError {
	return httpResponseError{ServiceName, err.Error()}
}

func newHTTPResponseOK() httpResponseOK {
	return httpResponseOK{ServiceName, true}
}

func writeJson(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, newHTTPResponseOK(), 200)
	}
}

func registerHandler(register registerable) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var service ServiceInfo
		values := r.URL.Query()

		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATH" {
			err := r.ParseForm()
			if err != nil {
				writeJson(w, newHTTPResponseError(err), 400)
				return
			}

			for name, value := range r.Form {
				if len(value) > 0 {
					values.Add(name, value[0])
				}
			}
		}

		service.Name = values.Get("service")
		service.Endpoint = values.Get("endpoint")
		service.HealthCheckTime = time.Now().Unix()
		service.Value = values.Get("value")

		if service.Name == "" {
			writeJson(w, newHTTPResponseError(errors.New("service name not provided")), 400)
			return
		}

		if service.Endpoint == "" {
			service.Endpoint = getRealIP(r)
		}

		register.RegisterService(service)

		writeJson(w, newHTTPResponseOK(), 200)
	}
}

func servicesHandler(services []ServiceInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, services, 200)
	}
}
