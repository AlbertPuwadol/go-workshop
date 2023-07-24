package adapter

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBConnection(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}

type IMongoAdapter interface {
	Find(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOptions) error
}

type mongodb struct {
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
}

func NewMongoDBAdapter(mongoClient *mongo.Client, mongoCollection *mongo.Collection) *mongodb {
	return &mongodb{mongoClient, mongoCollection}
}

func (m *mongodb) Find(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOptions) error {
	cursor, err := m.mongoCollection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, result)
	if err != nil {
		return err
	}
	return nil
}
