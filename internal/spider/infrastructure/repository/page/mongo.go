package pagerepository

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	loggerinterface "github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	config     spiderconfiginterface.Configurator
	logger     loggerinterface.Logger
	collection *mongo.Collection
}

func NewMongo(
	config spiderconfiginterface.Configurator,
	logger loggerinterface.Logger,
	mongodb *mongo.Database,
) *Mongo {
	return &Mongo{
		config:     config,
		logger:     logger,
		collection: mongodb.Collection(config.GetMongoPagesCollection()),
	}
}

func (r *Mongo) FindByURL(ctx context.Context, url string) (page *entity.Page, found bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.GetMongoRequestTimeout())
	defer cancel()

	filter := bson.M{"url": url}

	page = &entity.Page{}
	if err = r.collection.FindOne(ctx, filter).Decode(page); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, false, nil
		}
		return nil, false, r.logger.Error(ctx, err, nil)
	}

	return page, true, nil
}

func (r *Mongo) Save(ctx context.Context, page *entity.Page) error {
	ctx, cancel := context.WithTimeout(ctx, r.config.GetMongoRequestTimeout())
	defer cancel()

	res, err := r.collection.InsertOne(ctx, page, options.InsertOne())
	if err != nil {
		return r.logger.Error(ctx, err, nil)
	}

	if _, ok := res.InsertedID.(primitive.ObjectID); !ok {
		return r.logger.Error(ctx, err, nil)
	}

	return nil
}
