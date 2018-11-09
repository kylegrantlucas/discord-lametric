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
	Model Model `json:"model"`
}

type Model struct {
	Frames []Frame `json:"frames"`
}

type Frame struct {
	Icon int    `json:"icon"`
	Text string `json:"text"`
}

func (c Client) Notify(iconID int, message string) error {
	notification := Notification{Model: Model{Frames: []Frame{{Icon: iconID, Text: message}}}}

	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	client := http.Client{}
	request, err := http.NewRequest("POST", fmt.Sprintf("http://%v/api/v2/device/notifications", c.Host), bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	b64Key := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("dev:%v", c.APIKey)))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %v", b64Key))
	request.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}
	respData, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(respData))
	return nil
}
