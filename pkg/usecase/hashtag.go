package usecase

import (
	"context"

	"github.com/AlbertPuwadol/go-workshop/pkg/entity"
	"github.com/AlbertPuwadol/go-workshop/pkg/repository"
)

type IHashtag interface {
	GetAll(ctx context.Context) ([]entity.Hashtag, error)
	CreateQueue(queuename string) error
	Publish(queuename string, ctx context.Context, message string) error
}

type hashtag struct {
	repo repository.IHashtag
}

func NewHashtag(repo repository.IHashtag) *hashtag {
	return &hashtag{repo}
}

func (h hashtag) GetAll(ctx context.Context) ([]entity.Hashtag, error) {
	users, err := h.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (h hashtag) CreateQueue(queuename string) error {
	return h.repo.CreateQueue(queuename)
}

func (h hashtag) Publish(queuename string, ctx context.Context, message string) error {
	return h.repo.Publish(queuename, ctx, message)
}
