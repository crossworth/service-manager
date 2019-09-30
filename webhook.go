package servicemanager

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type webhookInfo struct {
	Old  []ServiceInfo `json:"old"`
	New  []ServiceInfo `json:"new"`
	Time int64         `json:"time"`
}

func notifyChanges(webhookUrls []*url.URL, maxWebhookTries int, httpClient *http.Client, old []ServiceInfo, new []ServiceInfo) {
	data := webhookInfo{
		Old:  old,
		New:  new,
		Time: time.Now().Unix(),
	}
	buffer, err := json.Marshal(&data)
	if err != nil {
		log.Printf("notifyChanges: error encoding webhook data %s\n", err)
		return
	}

	for _, u := range webhookUrls {
		for i := 0; i < maxWebhookTries; i++ {
			reader := bytes.NewBuffer(buffer)
			resp, err := httpClient.Post(u.String(), "application/json", reader)
			if err != nil {
				log.Printf("notifyChanges: httpClient.Post error %s\n", err)
				time.Sleep(5 * time.Second)
				continue
			}
			_ = resp.Body.Close()
			break
		}
	}
}
