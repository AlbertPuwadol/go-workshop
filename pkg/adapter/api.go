package adapter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AlbertPuwadol/go-workshop/pkg/entity"
)

type IApi interface {
	GetTimeline(hashtag string, pagesize int, cursor *string) (*entity.Timeline, error)
	GetInfo(userid string) (*entity.Info, error)
}

type api struct {
	url string
}

func NewApi(url string) *api {
	return &api{url}
}

func (a *api) GetTimeline(hashtag string, pagesize int, cursor *string) (*entity.Timeline, error) {
	method := "GET"

	client := &http.Client{}

	url := fmt.Sprintf("%s/thread/?hashtag=%s", a.url, hashtag)

	if pagesize != 0 {
		url += fmt.Sprintf("&page_size=%d", pagesize)
	}
	if cursor != nil {
		url += fmt.Sprintf("&cursor=%s", *cursor)
	}

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	timeline := &entity.Timeline{}

	err = json.Unmarshal(body, timeline)
	if err != nil {
		return nil, err
	}

	return timeline, nil
}

func (a *api) GetInfo(userid string) (*entity.Info, error) {
	method := "GET"

	client := &http.Client{}

	url := fmt.Sprintf("%s/account/%s", a.url, userid)

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	info := &entity.Info{}

	err = json.Unmarshal(body, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}
