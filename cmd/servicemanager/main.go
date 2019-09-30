package main

import (
	"github.com/crossworth/service-manager"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("using default http port %q\n", port)
	}

	webhook := os.Getenv("CHANGES_WEBHOOK")
	var webhookUrls []string
	if webhook != "" {
		webhookUrls = strings.Split(webhook, ",")
	}

	sm, err := servicemanager.New(time.Second*20, webhookUrls)
	if err != nil {
		log.Fatalf("error creating service manager  %s\n", err)
	}

	s := &http.Server{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + port,
		Handler:      sm,
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("error listing on port %q %s\n", port, err)
	}
}
