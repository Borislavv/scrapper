package spider

import (
	"context"
	"fmt"
	"github.com/Borislavv/scrapper/internal/spider/app/config"
	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	runnerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/job/runner/interface"
	schedulerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/job/scheduler/interface"
	pagecomparator "github.com/Borislavv/scrapper/internal/spider/domain/service/page/comparator"
	pagefinder "github.com/Borislavv/scrapper/internal/spider/domain/service/page/finder"
	pagesaver "github.com/Borislavv/scrapper/internal/spider/domain/service/page/saver"
	pagerepository "github.com/Borislavv/scrapper/internal/spider/infrastructure/repository/page"
	jobrunner "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/job/runner"
	jobscheduler "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/job/scheduler"
	"github.com/Borislavv/scrapper/internal/spider/infrastructure/service/page/scrapper"
	taskparser "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/task/parser"
	taskprovider "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/task/provider"
	taskrunner "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/task/runner"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
)

type Spider struct {
	ctx          context.Context
	config       spiderinterface.Config
	jobRunner    runnerinterface.JobRunner
	jobScheduler schedulerinterface.JobScheduler
}

func New(ctx context.Context) *Spider {
	cfg, err := spiderconfig.Load()
	if err != nil {
		panic(err)
	}

	// infrastructure
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf(
			"mongodb://%s:%s@%s:%d",
			cfg.GetMongoLogin(),
			cfg.GetMongoPassword(),
			cfg.GetMongoHost(),
			cfg.GetMongoPort(),
		),
	)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	go func() {
		<-ctx.Done()
		_ = mongoClient.Disconnect(ctx)
	}()

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	database := mongoClient.Database(cfg.GetMongoDatabase())

	// page dependencies
	pageRepo := pagerepository.New(cfg, database)
	pageSaver := pagesaver.New(pageRepo)
	pageFinder := pagefinder.New(pageRepo)
	pageScrapper := pagescrapper.New(cfg)
	pageComparator := pagecomparator.New()
	// task dependencies
	taskParser := taskparser.New(cfg)
	taskProvider := taskprovider.New(cfg, taskParser)
	taskRunner := taskrunner.New(pageSaver, pageFinder, pageScrapper, pageComparator)
	// job dependencies
	jobRunner := jobrunner.New(cfg, taskRunner, taskProvider)
	jobScheduler := jobscheduler.New(cfg)

	return &Spider{
		ctx:          ctx,
		config:       cfg,
		jobRunner:    jobRunner,
		jobScheduler: jobScheduler,
	}
}

func (s *Spider) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	for range s.jobScheduler.Manage(s.ctx) {
		s.jobRunner.Run(s.ctx)
	}
}
