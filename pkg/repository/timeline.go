package repository

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/AlbertPuwadol/go-workshop/pkg/adapter"
	"github.com/AlbertPuwadol/go-workshop/pkg/entity"
)

type ITimeline interface {
	CreateQueue(queuename string) error
	Publish(queuename string, ctx context.Context, data entity.Data) error
	GetTimeline(hashtag string, pagesize int, cursor *string) (*entity.Timeline, error)
	GetInfo(data *entity.Data) error
}

type timeline struct {
	apiAdapter      adapter.IApi
	rabbitMQAdapter adapter.IRabbitMQ
	accountMapper   map[string]entity.Info
	mutex           *sync.Mutex
}

func NewTimeline(apiAdapter adapter.IApi, rabbitMQAdapter adapter.IRabbitMQ) *timeline {
	accountMapper := make(map[string]entity.Info)
	var mutex sync.Mutex
	return &timeline{apiAdapter, rabbitMQAdapter, accountMapper, &mutex}
}

func (t *timeline) CreateQueue(queuename string) error {
	return t.rabbitMQAdapter.CreateQueue(queuename)
}

func (t *timeline) GetTimeline(hashtag string, pagesize int, cursor *string) (*entity.Timeline, error) {
	return t.apiAdapter.GetTimeline(hashtag, pagesize, cursor)
}

func (t *timeline) GetInfo(data *entity.Data) error {
	if v, ok := t.accountMapper[data.UserID]; ok {
		data.User = v
	} else {
		info, err := t.apiAdapter.GetInfo(data.UserID)
		if err != nil {
			return err
		}
		t.mutex.Lock()
		t.accountMapper[data.UserID] = *info
		t.mutex.Unlock()
		data.User = *info
	}
	return nil
}

func (t *timeline) Publish(queuename string, ctx context.Context, data entity.Data) error {
	job, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return t.rabbitMQAdapter.Publish(queuename, ctx, job)
}
