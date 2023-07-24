package main

import (
	"context"
	"log"
	"time"

	"github.com/AlbertPuwadol/go-workshop/config"
	"github.com/AlbertPuwadol/go-workshop/pkg/adapter"
	"github.com/AlbertPuwadol/go-workshop/pkg/repository"
	"github.com/AlbertPuwadol/go-workshop/pkg/usecase"
)

func main() {
	cfg := config.NewConfig()

	log.Println(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rabbitch, err := adapter.NewRabbitMQChannel(ctx, cfg.RabbitMQURI)
	if err != nil {
		log.Panic(err)
	}
	defer rabbitch.Close()

	rabbitMQAdapter := adapter.NewRabbitMQAdapter(rabbitch)

	mongodbClient, err := adapter.NewMongoDBConnection(ctx, cfg.MongoDBURI)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err = mongodbClient.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	hashtagCollection := mongodbClient.Database("go_workshop").Collection("hashtag")
	mongoDBAdapter := adapter.NewMongoDBAdapter(mongodbClient, hashtagCollection)

	repository := repository.Newhashtag(mongoDBAdapter, rabbitMQAdapter)

	usecase := usecase.NewHashtag(repository)

	usecase.CreateQueue(cfg.IntervalQueue)

	result, err := usecase.GetAll(ctx)
	if err != nil {
		log.Println(err)
	}

	log.Print(result)

	for _, v := range result {
		usecase.Publish(cfg.IntervalQueue, ctx, v.Keyword)
		log.Printf("Publish to queue: %s message: %s\n", cfg.IntervalQueue, v.Keyword)
	}

}
