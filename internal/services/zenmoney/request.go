package zenmoney

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type Diff struct {
	CurrentClientTimestamp int64 `json:"currentClientTimestamp"`
	ServerTimestamp        int64 `json:"serverTimestamp"`
}

func (api *Api) Diff() (*http.Response, error) {
	token, err := api.Init()
	if err != nil {
		return nil, err
	}

	bearer := "Bearer " + token

	d := Diff{time.Now().Unix(), 0}

	diff, _ := json.Marshal(d)

	log.Print(d)
	log.Print(string(diff))

	req, errorReq := http.NewRequest("POST", BASE_URL, bytes.NewReader(diff))

	if errorReq != nil {
		return nil, errorReq
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, errorResp := client.Do(req)

	if errorResp != nil {
		return nil, errorResp
	}

	if resp.StatusCode != 200 {
		log.Print(resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		log.Print(string(body))
		return nil, errors.New(resp.Status)
	}

	return resp, errorResp
}
