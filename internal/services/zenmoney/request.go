package zenmoney

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Diff struct {
	CurrentClientTimestamp int64 `json:"currentClientTimestamp"`
	ServerTimestamp        int64 `json:"serverTimestamp"`
}

func (api *Api) Diff() (*Response, error) {
	startTime := time.Now()
	log.Printf("[Diff] Начало запроса")

	token, err := api.Init()
	if err != nil {
		log.Printf("[Diff] Init заняло: %v, ошибка: %v", time.Since(startTime), err)
		return nil, err
	}
	log.Printf("[Diff] Init заняло: %v", time.Since(startTime))

	bearer := "Bearer " + token

	d := Diff{time.Now().Unix(), 0}
	diff, _ := json.Marshal(d)

	req, errorReq := http.NewRequest("POST", BASE_URL, bytes.NewReader(diff))
	if errorReq != nil {
		log.Printf("[Diff] Создание запроса заняло: %v, ошибка: %v", time.Since(startTime), errorReq)
		return nil, errorReq
	}
	log.Printf("[Diff] Создание запроса заняло: %v", time.Since(startTime))

	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	log.Printf("[Diff] Отправка запроса:")
	resp, errorResp := client.Do(req)
	log.Printf("[Diff] Отправка HTTP запроса заняло: %v", time.Since(startTime))

	if errorResp != nil {
		log.Printf("[Diff] Ошибка выполнения запроса: %v, всего заняло: %v", errorResp, time.Since(startTime))
		return nil, errorResp
	}

	if resp.StatusCode != 200 {
		log.Printf("[Diff] Нестатус 200: %v, всего заняло: %v", resp.StatusCode, time.Since(startTime))
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[Diff] Тело ошибки: %s", string(body))
		log.Printf("[Diff] Чтение тела ошибки заняло: %v", time.Since(startTime))
		return nil, errors.New(resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	log.Printf("[Diff] Чтение ответа заняло: %v", time.Since(startTime))

	var result Response

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[Diff] Ошибка парсинга JSON: %v, всего заняло: %v", err, time.Since(startTime))
		fmt.Println("Can not unmarshal JSON")
	}

	log.Printf("[Diff] Запрос Diff завершен, всего заняло: %v", time.Since(startTime))
	return &result, errorResp
}
