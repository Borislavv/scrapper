package pagerepository

import (
	"context"
	"errors"
	"gitlab.xbet.lan/web-backend/php/spider/internal/shared/domain/entity"
	spiderconfiginterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/app/config/interface"
	pagerepositoryinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/repository/interface"
	"gitlab.xbet.lan/web-backend/php/spider/internal/spider/infrastructure/logger/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	config     spiderconfiginterface.Configurator
	logger     logger.Logger
	collection *mongo.Collection
}

func NewMongo(
	config spiderconfiginterface.Configurator,
	logger logger.Logger,
	mongodb *mongo.Database,
) *Mongo {
	return &Mongo{
		config:     config,
		logger:     logger,
		collection: mongodb.Collection(config.GetMongoPagesCollection()),
	}
}

func (r *Mongo) FindOneLatestByURL(ctx context.Context, url string) (page *entity.Page, found bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.GetMongoRequestTimeout())
	defer cancel()

	filter := bson.M{"url": url}
	opts := options.FindOne().SetSort(bson.D{{"version", -1}})

	page = &entity.Page{}
	if err = r.collection.FindOne(ctx, filter, opts).Decode(page); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, false, nil
		}
		return nil, false, r.logger.Error(ctx, pagerepositoryinterface.FindByURLError, logger.Fields{
			"url": url,
			"err": err.Error(),
		})
	}

	return page, true, nil
}

func (r *Mongo) Save(ctx context.Context, page *entity.Page) error {
	ctx, cancel := context.WithTimeout(ctx, r.config.GetMongoRequestTimeout())
	defer cancel()

	res, err := r.collection.InsertOne(ctx, page, options.InsertOne())
	if err != nil {
		return r.logger.Error(ctx, pagerepositoryinterface.SaveError, logger.Fields{
			"url":       page.URL,
			"userAgent": page.UserAgent,
			"err":       err.Error(),
		})
	}

	if _, ok := res.InsertedID.(primitive.ObjectID); !ok {
		return r.logger.Error(ctx, pagerepositoryinterface.InsertedIDCastError, logger.Fields{
			"url":        page.URL,
			"userAgent":  page.UserAgent,
			"version":    page.Version,
			"insertedId": res.InsertedID,
		})
	}

	return nil
}

func (r *Mongo) Update(ctx context.Context, page *entity.Page) error {
	ctx, cancel := context.WithTimeout(ctx, r.config.GetMongoRequestTimeout())
	defer cancel()

	page.UpdatedAt = time.Now()

	updateData, err := bson.Marshal(page)
	if err != nil {
		return r.logger.Error(ctx, pagerepositoryinterface.UpdateMarshalError, logger.Fields{
			"url":       page.URL,
			"userAgent": page.UserAgent,
			"version":   page.Version,
			"err":       err.Error(),
		})
	}

	var update bson.M
	if err = bson.Unmarshal(updateData, &update); err != nil {
		return r.logger.Error(ctx, pagerepositoryinterface.UpdateUnmarshalError, logger.Fields{
			"url":       page.URL,
			"userAgent": page.UserAgent,
			"version":   page.Version,
			"err":       err.Error(),
		})
	}

	update = bson.M{"$set": update}

	_, err = r.collection.UpdateByID(ctx, page.GetID().ID, update, options.Update())
	if err != nil {
		return r.logger.Error(ctx, pagerepositoryinterface.UpdateError, logger.Fields{
			"url":       page.URL,
			"userAgent": page.UserAgent,
			"version":   page.Version,
			"err":       err.Error(),
		})
	}

	return nil
}
