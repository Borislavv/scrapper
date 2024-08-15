package spider

import (
	"context"
	sharedconfig "github.com/Borislavv/scrapper/internal/shared/app/config"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/database"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/logger"
	"github.com/Borislavv/scrapper/internal/spider/app/config"
	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	jobrunner "github.com/Borislavv/scrapper/internal/spider/domain/service/job/runner"
	runnerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/job/runner/interface"
	schedulerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/job/scheduler/interface"
	pagecomparator "github.com/Borislavv/scrapper/internal/spider/domain/service/page/comparator"
	pageconsumer "github.com/Borislavv/scrapper/internal/spider/domain/service/page/consumer"
	pageprovider "github.com/Borislavv/scrapper/internal/spider/domain/service/page/provider"
	taskconsumer "github.com/Borislavv/scrapper/internal/spider/domain/service/task/consumer"
	taskprovider "github.com/Borislavv/scrapper/internal/spider/domain/service/task/provider"
	taskrunner "github.com/Borislavv/scrapper/internal/spider/domain/service/task/runner"
	pagerepository "github.com/Borislavv/scrapper/internal/spider/infrastructure/repository/page"
	jobscheduler "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/job/scheduler"
	pageparser "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/page/parser"
	pagescanner "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/page/scanner"
	taskparser "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/task/parser"
	"log"
	"sync"
)

type Spider struct {
	ctx          context.Context
	config       spiderinterface.Configurator
	jobRunner    runnerinterface.JobRunner
	jobScheduler schedulerinterface.JobScheduler
}

func New(ctx context.Context) (*Spider, error) {
	sharedCfg, err := sharedconfig.Load()
	if err != nil {
		return nil, err
	}

	logrus, cancel, err := logger.NewLogrus(sharedCfg)
	if err != nil {
		return nil, err
	}
	defer cancel()

	spiderCfg, err := spiderconfig.Load()
	if err != nil {
		return nil, logrus.Fatal(ctx, err, nil)
	}

	// infrastructure
	mongodb, err := database.NewMongo(logrus).Connect(ctx, sharedCfg)
	if err != nil {
		return nil, logrus.Fatal(ctx, err, nil)
	}

	// page infra. dependencies
	pageRepo := pagerepository.NewMongo(spiderCfg, logrus, mongodb)
	pageParser := pageparser.NewHTML(logrus)
	pageScanner := pagescanner.NewHTTP(spiderCfg, pageParser, logrus)

	// page domain dependencies
	pageComparator := pagecomparator.NewEqual()
	pageConsumer := pageconsumer.NewParallel(logrus, pageRepo, pageComparator)
	pageProvider := pageprovider.NewChan(spiderCfg, pageScanner)

	// task dependencies
	taskParser := taskparser.NewCSV(spiderCfg, logrus)
	taskProvider, err := taskprovider.NewParallel(ctx, logrus, spiderCfg, taskParser)
	if err != nil {
		return nil, logrus.Fatal(ctx, err, nil)
	}
	taskRunner := taskrunner.New(spiderCfg, pageProvider, pageConsumer)
	taskConsumer := taskconsumer.NewParallel(taskRunner)

	// job dependencies
	jobRunner := jobrunner.New(spiderCfg, taskConsumer, taskProvider)
	jobScheduler := jobscheduler.NewTicker(spiderCfg)

	return &Spider{
		ctx:          ctx,
		config:       spiderCfg,
		jobRunner:    jobRunner,
		jobScheduler: jobScheduler,
	}, nil
}

func (s *Spider) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("starting a new job")
	s.jobRunner.Run(s.ctx)
	log.Println("finished a job")

	//for range s.jobScheduler.Manage(s.ctx) {
	//	log.Println("starting a new job")
	//	s.jobRunner.Run(s.ctx)
	//	log.Println("finished a job")
	//}
}
