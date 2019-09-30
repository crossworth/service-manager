package client

import (
	"bytes"
	"encoding/json"
	servicemanager "github.com/crossworth/service-manager"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	defaultHttpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	defaultLogger = log.New(os.Stdout, "servicemanager", log.LstdFlags)
)

// Options struct
type Options struct {
	Name            string
	ManagerEndpoint string
	UpdateInterval  time.Duration
	Endpoint        string
	Value           func() string
	HttpClient      *http.Client
	Logger          *log.Logger
}

// Register a new service
func Register(options Options) error {
	if options.HttpClient == nil {
		options.HttpClient = defaultHttpClient
	}

	if options.Logger == nil {
		options.Logger = defaultLogger
	}

	if options.Name == "" {
		return errors.New("servicemanager name not set")
	}

	if options.ManagerEndpoint == "" {
		return errors.New("servicemanager manager endpoint not set")
	}

	if options.UpdateInterval == 0 {
		options.UpdateInterval = 20 * time.Second
	}

	if options.Endpoint == "" {
		var err error
		options.Endpoint, err = GetLocalIP()
		if err != nil {
			return err
		}
	}

	if options.Value == nil {
		options.Value = func() string {
			return ""
		}
	}

	err := checkServiceManager(options.HttpClient, options.ManagerEndpoint)

	if err != nil {
		return errors.Wrap(err, "servicemanager endpoint error")
	}

	notifyServiceManager(options)

	go (func() {
		for {
			notifyServiceManager(options)
			time.Sleep(options.UpdateInterval)
		}
	})()

	return nil
}

func checkServiceManager(httpClient *http.Client, managerEndpoint string) error {
	response, err := httpClient.Get(managerEndpoint)

	if err != nil {
		return err
	}

	var httpResponse map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&httpResponse)

	if err != nil {
		return err
	}

	if httpResponse["service"] != servicemanager.ServiceName {
		return errors.New("service-manager signature mismatch")
	}

	return nil
}

func notifyServiceManager(options Options) {
	url := options.ManagerEndpoint + "/register?service=" + options.Name + "&endpoint=" + options.Endpoint

	response, err := options.HttpClient.Post(url, "application/json", bytes.NewReader([]byte(options.Value())))
	if err != nil {
		options.Logger.Printf("notifyServiceManager: error netClient.Post %s\n", err)
		return
	}
	_ = response.Body.Close()
}

// GetLocalIP returns the preferred outbound ip address
func GetLocalIP() (string, error) {
	// source: https://stackoverflow.com/a/37382208
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return "", errors.Wrap(err, "could not get local ip")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
