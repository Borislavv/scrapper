package database

import (
	"context"
	"fmt"
	logger "github.com/Borislavv/scrapper/internal/shared/domain/service/logger/interface"
	databaseconfiginterface "github.com/Borislavv/scrapper/internal/shared/infrastructure/database/config/interface"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	logger logger.Logger
}

func NewMongo(logger logger.Logger) *Mongo {
	return &Mongo{logger: logger}
}

func (m *Mongo) Connect(ctx context.Context, cfg databaseconfiginterface.Configurator) (*mongo.Database, error) {
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
		m.logger.FatalMsg(ctx, "mongodb connection failed", logger.Fields{
			"err": err.Error(),
		})
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = mongoClient.Disconnect(ctx)
	}()

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		m.logger.FatalMsg(ctx, "mongodb ping failed", logger.Fields{
			"err": err.Error(),
		})
		return nil, err
	}

	return mongoClient.Database(cfg.GetMongoDatabase()), nil
}
