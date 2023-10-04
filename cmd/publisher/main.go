package main

import (
	"context"
	"log"
	"time"

	"github.com/AlbertPuwadol/go-workshop/config"
	"github.com/AlbertPuwadol/go-workshop/pkg/adapter"
	"github.com/AlbertPuwadol/go-workshop/pkg/repository"
	"github.com/AlbertPuwadol/go-workshop/pkg/usecase"
	"github.com/AlbertPuwadol/test-go-util/mongo"
)

func main() {
	log.Println(mongo.Hello())
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

	repository := repository.NewHashtag(mongoDBAdapter, rabbitMQAdapter)

	usecase := usecase.NewHashtag(repository)

	usecase.CreateQueue(cfg.IntervalQueue)

	result, err := usecase.GetAll(ctx)
	if err != nil {
		log.Println(err)
	}

	log.Println(result)

	usecase.Publish(cfg.IntervalQueue, ctx, result)
	log.Printf("Publish to queue: %s success\n", cfg.IntervalQueue)

}
