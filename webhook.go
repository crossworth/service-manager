package servicemanager

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
)

type webhookInfo struct {
	Old  []ServiceInfo `json:"old"`
	New  []ServiceInfo `json:"new"`
	Time int64         `json:"time"`
}

func (s *Server) notifyChanges(old []ServiceInfo, new []ServiceInfo) {
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

	reader := bytes.NewReader(buffer)

	for _, u := range s.webhookUrls {
		for i := 0; i <= s.maxWebhookTries; i++ {
			resp, err := s.httpClient.Post(u.String(), "application/json", reader)
			if err != nil {
				log.Printf("notifyChanges: httpClient.Post error %s\n", err)
				time.Sleep(5 * time.Second)
				continue
			}
			defer resp.Body.Close()
		}
	}
}