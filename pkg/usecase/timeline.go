package usecase

import (
	"context"
	"log"
	"sync"

	"github.com/AlbertPuwadol/go-workshop/pkg/entity"
	"github.com/AlbertPuwadol/go-workshop/pkg/repository"
)

type ITimeline interface {
	CreateQueue(queuename entity.Hashtag) error
	Publish(queuename string, ctx context.Context, data []entity.Data) error
	GetTimeline(hashtag string) ([]entity.Data, error)
}

type timeline struct {
	repo repository.ITimeline
}

func NewTimeline(repo repository.ITimeline) *timeline {
	return &timeline{repo}
}

func (t timeline) CreateQueue(queuename string) error {
	return t.repo.CreateQueue(queuename)
}

func (t timeline) GetTimeline(hashtag entity.Hashtag) ([]entity.Data, error) {
	pagesize := 100
	var cursor *string
	var result []entity.Data
	var wg sync.WaitGroup
	for {
		res, err := t.repo.GetTimeline(hashtag.Keyword, pagesize, cursor)
		if err != nil {
			return nil, err
		}

		wg.Add(len(res.Data))
		channel := make(chan entity.Data, len(res.Data))
		for _, d := range res.Data {
			go func(data entity.Data) {
				defer wg.Done()
				err := t.repo.GetInfo(&data)
				if err != nil {
					log.Println(err)
					return
				}
				channel <- data
			}(d)
		}

		go func() {
			wg.Wait()
			close(channel)
		}()

		for item := range channel {
			result = append(result, item)
		}

		cursor = res.NextPage

		if cursor == nil {
			break
		}
		log.Printf("Hashtag: %s  has next cursor\n", hashtag.Keyword)
	}
	return result, nil
}

func (t timeline) Publish(queuename string, ctx context.Context, data []entity.Data) error {
	for _, d := range data {
		if err := t.repo.Publish(queuename, ctx, d); err != nil {
			return err
		}
	}
	return nil
}
