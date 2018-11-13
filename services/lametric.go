package lametric

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	Host   string
	APIKey string
}

type Notification struct {
	IconType *string `json:"icon_type,omitempty"`
	Priority *string `json:"priority,omitempty"`
	Model    Model   `json:"model"`
}

type Model struct {
	Frames []Frame `json:"frames,omitempty"`
	Sound  *Sound  `json:"sound,omitempty"`
}

type Frame struct {
	Icon *string `json:"icon"`
	Text string  `json:"text"`
}

type Sound struct {
	Category string `json:"category,omitempty"`
	ID       string `json:"id,omitempty"`
	Repeat   int    `json:"repeat,omitempty"`
}

func (c Client) Notify(notification Notification) error {
	log.Print(notification)
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	client := http.Client{}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%v:8080/api/v2/device/notifications", c.Host),
		bytes.NewReader(body),
	)
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	b64Key := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("dev:%v", c.APIKey)))

	request.Header.Add("Authorization", fmt.Sprintf("Basic %v", b64Key))
	request.Header.Add("Content-Type", "application/json")
	log.Print(request)
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}
	respData, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(respData))
	return nil
}
