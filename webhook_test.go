package servicemanager

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func asUrlSlice(t *testing.T, urlString string) []*url.URL {
	t.Helper()

	u, err := url.ParseRequestURI(urlString)
	if err != nil {
		t.Fatal(err)
	}

	return []*url.URL{u}
}

func readBody(t *testing.T, r *http.Request) string {
	t.Helper()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	return string(bytes)
}

func TestServer_notifyChanges(t *testing.T) {
	maxTries := 1
	httpClient := http.DefaultClient

	t.Run("Empty services", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)

			body := readBody(t, r)
			if !strings.Contains(body, `{"old":[],"new":[],"time":`) {
				t.Errorf("got wrong webhook payload %q", body)
			}
		}))
		defer testServer.Close()

		notifyChanges(asUrlSlice(t, testServer.URL), maxTries, httpClient, []ServiceInfo{}, []ServiceInfo{})
	})

	t.Run("New services", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)

			body := readBody(t, r)
			if !strings.Contains(body, `{"old":[],"new":[{"name":"Test","endpoint":"http://10.0.0.10:3030","value":"this-is-a-value"}],"time":`) {
				t.Errorf("got wrong webhook payload %q", body)
			}
		}))
		defer testServer.Close()

		notifyChanges(asUrlSlice(t, testServer.URL), maxTries, httpClient, []ServiceInfo{}, []ServiceInfo{
			{
				Name:     "Test",
				Endpoint: "http://10.0.0.10:3030",
				Value:    "this-is-a-value",
			},
		})
	})
}
