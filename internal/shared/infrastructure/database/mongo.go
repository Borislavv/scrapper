package database

import (
	"context"
	"fmt"

	sharedconfiginterface "github.com/Borislavv/scrapper/internal/shared/app/config/interface"
	loggerinterface "github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	logger loggerinterface.Logger
}

func NewMongo(logger loggerinterface.Logger) *Mongo {
	return &Mongo{logger: logger}
}

func (m *Mongo) Connect(ctx context.Context, cfg sharedconfiginterface.Configurator) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/?authSource=%s",
		cfg.GetMongoLogin(),
		cfg.GetMongoPassword(),
		cfg.GetMongoHost(),
		cfg.GetMongoPort(),
		cfg.GetMongoDatabase(),
	))

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		m.logger.FatalMsg(ctx, "mongodb connection failed", loggerinterface.Fields{
			"err": err.Error(),
		})
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = mongoClient.Disconnect(ctx)
	}()

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		m.logger.FatalMsg(ctx, "mongodb ping failed", loggerinterface.Fields{
			"err": err.Error(),
		})
		return nil, err
	}

	return mongoClient.Database(cfg.GetMongoDatabase()), nil
}
