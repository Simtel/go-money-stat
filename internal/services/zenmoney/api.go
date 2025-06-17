package zenmoney

import (
	"errors"
	"money-stat/internal/config"
	"net/http"
)

var BASE_URL string = "https://api.zenmoney.app/v8/diff/"

var Token string

type Api struct {
	client *http.Client
}

type ApiInterface interface {
	Diff() (*Response, error)
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
