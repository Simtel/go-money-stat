package zenmoney

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Diff struct {
	CurrentClientTimestamp int64 `json:"currentClientTimestamp"`
	ServerTimestamp        int64 `json:"serverTimestamp"`
}

type Response struct {
	Account []struct {
		Id         string  `json:"id"`
		Title      string  `json:"title"`
		Balance    float64 `json:"balance"`
		Instrument int64   `json:"instrument"`
	} `json:"Account"`
}

func (api *Api) Diff() (*http.Response, error) {
	token, err := api.Init()
	if err != nil {
		return nil, err
	}

	bearer := "Bearer " + token

	d := Diff{time.Now().Unix(), 0}
	diff, _ := json.Marshal(d)

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

	body, err := io.ReadAll(resp.Body) // response body is []byte

	var result Response

	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	tableData := pterm.TableData{
		{"Счет", "Баланс", "Валюта"},
		{" ", " ", " "},
	}

	for _, account := range result.Account {
		tableData = append(tableData, []string{account.Title, strconv.FormatFloat(account.Balance, 'f', 2, 64), strconv.FormatInt(account.Instrument, 16)})
	}

	pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()

	return resp, errorResp
}
