package pagerepository

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	pagerepositoryinterface "github.com/Borislavv/scrapper/internal/spider/domain/repository/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/url"
)

type PageRepository struct {
	config     spiderconfiginterface.Config
	collection *mongo.Collection
}

func New(config spiderconfiginterface.Config, mongodb *mongo.Database) *PageRepository {
	return &PageRepository{
		config:     config,
		collection: mongodb.Collection(config.GetMongoPagesCollection()),
	}
}

func (r *PageRepository) FindByURL(ctx context.Context, url url.URL) (*entity.Page, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.GetMongoRequestTimeout())
	defer cancel()

	filter := bson.M{
		"url": url.String(),
	}

	page := &entity.Page{}
	if err := r.collection.FindOne(ctx, filter).Decode(page); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		log.Println("PageRepository: " + err.Error())
		return nil, err
	}

	return page, nil
}

func (r *PageRepository) Save(ctx context.Context, page *entity.Page) error {
	ctx, cancel := context.WithTimeout(ctx, r.config.GetMongoRequestTimeout())
	defer cancel()

	res, err := r.collection.InsertOne(ctx, page, options.InsertOne())
	if err != nil {
		log.Println("PageRepository: " + err.Error())
		return err
	}

	if _, ok := res.InsertedID.(primitive.ObjectID); !ok {
		log.Println("PageRepository: " + pagerepositoryinterface.InsertError.Error())
		return pagerepositoryinterface.InsertError
	}

	return nil
}
