package servicemanager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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
	w.Header().Set("Content-Type", "application/json")
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
		service.Value = values.Get("value")

		if service.Value == "" {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			service.Value = string(bodyBytes)
		}

		if service.Name == "" {
			writeJson(w, newHTTPResponseError(errors.New("service not provided")), 400)
			return
		}

		if service.Endpoint == "" {
			writeJson(w, newHTTPResponseError(errors.New("endpoint not provided")), 400)
			return
		}

		register.RegisterService(service)

		writeJson(w, newHTTPResponseOK(), 200)
	}
}

func servicesHandler(servicesFunc func() []ServiceInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, servicesFunc(), 200)
	}
}
