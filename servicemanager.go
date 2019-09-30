package servicemanager

import (
	"fmt"
	"github.com/carlescere/scheduler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/patrickmn/go-cache"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

// The service name and signature
const ServiceName = "service-manager"

// ServiceInfo is the representation of a registered service
type ServiceInfo struct {
	Name            string      `json:"name"`
	Endpoint        string      `json:"endpoint"`
	HealthCheckTime int64       `json:"health_check_time"`
	Value           interface{} `json:"value,omitempty"`
}

// Server is main application struct
type Server struct {
	httpClient        *http.Client
	services          *cache.Cache
	servicesLastCheck []ServiceInfo
	webhookUrls       []*url.URL
	maxWebhookTries   int
	chi.Router
}

type registerable interface {
	RegisterService(service ServiceInfo)
}

// New create a new service manager with the provided
// time to live and webhook urls, if any webhook url
// is provided it start a background job to check
// and notify services changes to the webhook urls
func New(timeToLive time.Duration, webhookUrls []string) (*Server, error) {
	s := &Server{
		maxWebhookTries: 3,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		services: cache.New(timeToLive, timeToLive*2),
		Router:   chi.NewRouter(),
	}

	for _, u := range webhookUrls {
		u, err := url.ParseRequestURI(u)
		if err != nil {
			return &Server{}, fmt.Errorf("invalid webhook url %q", err)
		}

		s.webhookUrls = append(s.webhookUrls, u)
	}

	if len(s.webhookUrls) > 0 {
		_, err := scheduler.Every(int(timeToLive.Seconds())).Seconds().Run(s.checkForServiceChanges)
		return s, err
	}

	s.setHttpHandlers()
	return s, nil
}

// SetHttpClient allows setting a custom http client
// used on the webhook calls
func (s *Server) SetHttpClient(client *http.Client) {
	s.httpClient = client
}

// GetServices returns all the registered services
func (s *Server) GetServices() []ServiceInfo {
	items := s.services.Items()

	r := make([]ServiceInfo, 0)

	for _, e := range items {
		r = append(r, e.Object.(ServiceInfo))
	}

	return r
}

// RegisterService register a new service or replace it with
// the default time to live
func (s *Server) RegisterService(service ServiceInfo) {
	_, exists := s.services.Get(service.Endpoint)
	if exists {
		_ = s.services.Replace(service.Endpoint, service, cache.DefaultExpiration)
		return
	}

	_ = s.services.Add(service.Endpoint, service, cache.DefaultExpiration)
}

func (s *Server) setHttpHandlers() {
	s.Router.Use(
		middleware.Recoverer,
		middleware.Logger,
		middleware.DefaultCompress,
	)
	s.Handle("/", homeHandler())
	s.Handle("/register", registerHandler(s))
	s.Handle("/services", servicesHandler(s.GetServices))
}

func (s *Server) checkForServiceChanges() {
	newServices := s.GetServices()
	oldServices := s.servicesLastCheck

	hasChanges := false

	if !reflect.DeepEqual(oldServices, newServices) {
		hasChanges = true
	}

	if hasChanges && len(s.webhookUrls) > 0 {
		// https://github.com/go101/go101/wiki/How-to-efficiently-clone-a-slice%3F
		go s.notifyChanges(append(oldServices[:0:0], oldServices...), append(newServices[:0:0], newServices...))
	}

	s.servicesLastCheck = newServices
}
