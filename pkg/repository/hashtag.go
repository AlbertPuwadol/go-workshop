package repository

import (
	"context"
	"encoding/json"

	"github.com/AlbertPuwadol/go-workshop/pkg/adapter"
	"github.com/AlbertPuwadol/go-workshop/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type IHashtag interface {
	GetAll(ctx context.Context) ([]entity.Hashtag, error)
	CreateQueue(queuename string) error
	Publish(queuename string, ctx context.Context, hashtag entity.Hashtag) error
}

type hashtag struct {
	mongoDBAdapter  adapter.IMongoAdapter
	rabbitMQAdapter adapter.IRabbitMQ
}

func NewHashtag(mongoDBAdapter adapter.IMongoAdapter, rabbitMQAdapter adapter.IRabbitMQ) *hashtag {
	return &hashtag{mongoDBAdapter, rabbitMQAdapter}
}

func (h hashtag) GetAll(ctx context.Context) ([]entity.Hashtag, error) {

	var users []entity.Hashtag

	err := h.mongoDBAdapter.Find(ctx, &users, bson.D{})

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (h hashtag) CreateQueue(queuename string) error {
	return h.rabbitMQAdapter.CreateQueue(queuename)
}

func (h hashtag) Publish(queuename string, ctx context.Context, hashtag entity.Hashtag) error {
	job, err := json.Marshal(hashtag)
	if err != nil {
		return err
	}
	return h.rabbitMQAdapter.Publish(queuename, ctx, job)
}
