package zenmoney

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"money-stat/internal/config"
	"net/http"
	"time"
)

var BASE_URL string = "https://api.zenmoney.ru/v8/diff/"

var Token string

type Api struct {
	client *http.Client
}

type ApiInterface interface {
	Diff() (*Response, error)
	DiffSince(timestamp int64) (*Response, error)
}

func NewApi(client *http.Client) ApiInterface {
	return &Api{
		client: client,
	}
}

func (request *Api) Init() (string, error) {
	conf := config.New()
	token := conf.ZenMoney.Token
	if token == "" {
		return "", errors.New("you need to set ZENMONEY TOKEN environment variable")
	}
	return token, nil
}

// DiffSince запрашивает изменения с указанного timestamp
func (api *Api) DiffSince(timestamp int64) (*Response, error) {
	startTime := time.Now()
	log.Printf("[DiffSince] Начало запроса с timestamp=%d", timestamp)

	token, err := api.Init()
	if err != nil {
		log.Printf("[DiffSince] Init заняло: %v, ошибка: %v", time.Since(startTime), err)
		return nil, err
	}
	log.Printf("[DiffSince] Init заняло: %v", time.Since(startTime))

	bearer := "Bearer " + token

	// currentClientTimestamp — текущее клиентское время (секунды)
	clientTimestamp := time.Now().Unix()
	log.Printf("[DiffSince] currentClientTimestamp=%d, serverTimestamp=%d (в секундах)", clientTimestamp, timestamp)

	d := Diff{CurrentClientTimestamp: clientTimestamp, ServerTimestamp: timestamp}
	diff, err := json.Marshal(d)
	if err != nil {
		log.Printf("[DiffSince] Ошибка маршалинга JSON: %v, всего заняло: %v", err, time.Since(startTime))
		return nil, fmt.Errorf("ошибка маршалинга JSON запроса DiffSince: %w", err)
	}

	req, errorReq := http.NewRequest("POST", BASE_URL, bytes.NewReader(diff))
	if errorReq != nil {
		log.Printf("[DiffSince] Создание запроса заняло: %v, ошибка: %v", time.Since(startTime), errorReq)
		return nil, errorReq
	}
	log.Printf("[DiffSince] Создание запроса заняло: %v", time.Since(startTime))

	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	log.Printf("[DiffSince] Отправка HTTP запроса:")
	resp, errorResp := api.client.Do(req)
	log.Printf("[DiffSince] Отправка HTTP запроса заняло: %v", time.Since(startTime))

	if errorResp != nil {
		log.Printf("[DiffSince] Ошибка выполнения запроса: %v, всего заняло: %v", errorResp, time.Since(startTime))
		return nil, errorResp
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("[DiffSince] Нестатус 200: %v, всего заняло: %v", resp.StatusCode, time.Since(startTime))
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("[DiffSince] Ошибка чтения тела ошибки: %v", readErr)
			return nil, errors.New(resp.Status)
		}
		log.Printf("[DiffSince] Тело ошибки: %s", string(body))
		log.Printf("[DiffSince] Чтение тела ошибки заняло: %v", time.Since(startTime))
		return nil, errors.New(resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("[DiffSince] Ошибка чтения ответа: %v, всего заняло: %v", readErr, time.Since(startTime))
		return nil, fmt.Errorf("ошибка чтения тела ответа DiffSince: %w", readErr)
	}
	log.Printf("[DiffSince] Чтение ответа заняло: %v", time.Since(startTime))

	var result Response

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[DiffSince] Ошибка парсинга JSON: %v, всего заняло: %v", err, time.Since(startTime))
		return nil, err
	}

	log.Printf("[DiffSince] Запрос DiffSince завершен, всего заняло: %v. Получено транзакций: %d", time.Since(startTime), len(result.Transaction))
	return &result, nil
}
