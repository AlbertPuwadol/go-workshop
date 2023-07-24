package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/AlbertPuwadol/go-workshop/config"
	"github.com/AlbertPuwadol/go-workshop/pkg/adapter"
	"github.com/AlbertPuwadol/go-workshop/pkg/entity"
	"github.com/AlbertPuwadol/go-workshop/pkg/repository"
	"github.com/AlbertPuwadol/go-workshop/pkg/usecase"
)

func main() {
	cfg := config.NewConfig()

	log.Println(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rabbitconn, rabbitch, err := adapter.NewRabbitMQChannel(ctx, cfg.RabbitMQURI)
	if err != nil {
		log.Panic(err)
	}

	rabbitMQAdapter := adapter.NewRabbitMQAdapter(rabbitconn, rabbitch)

	defer rabbitMQAdapter.Close()

	apiAdapter := adapter.NewApi(cfg.ApiURL)

	repository := repository.NewTimeline(apiAdapter, rabbitMQAdapter)

	usecase := usecase.NewTimeline(repository)

	usecase.CreateQueue(cfg.ResultQueue)

	rabbitMQAdapter.Consume(cfg.IntervalQueue, func(job []byte) {
		var hashtag entity.Hashtag
		err := json.Unmarshal(job, &hashtag)
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Got job %+v\n", hashtag)

		data, err := usecase.GetTimeline(hashtag)
		if err != nil {
			log.Panic(err)
		}

		err = usecase.Publish(cfg.ResultQueue, ctx, data)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Publish Thread Hashtag: %s Success\n", hashtag.Keyword)
	})
}
