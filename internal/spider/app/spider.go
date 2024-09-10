package spider

import (
	"context"
	loggerinterface "github.com/Borislavv/scrapper/internal/shared/domain/service/logger/interface"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/database"
	liveness "github.com/Borislavv/scrapper/internal/shared/infrastructure/liveness/interface"
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
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	ctx          context.Context
	cancel       context.CancelFunc
	config       spiderinterface.Configurator
	logger       loggerinterface.Logger
	loggerCancel loggerinterface.CancelFunc
	jobRunner    runnerinterface.JobRunner
	jobScheduler schedulerinterface.JobScheduler
	liveness     liveness.Prober
	mongo        *mongo.Database
}

func New(ctx context.Context, output loggerinterface.Outputer, liveness liveness.Prober) (*App, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel

	cfg, err := spiderconfig.Load()
	if err != nil {
		return nil, err
	}

	logrus, logrusCancel, err := logger.NewLogrus(cfg, output)
	if err != nil {
		return nil, err
	}

	// infrastructure
	mongodb, err := database.NewMongo(logrus).Connect(ctx, cfg)
	if err != nil {
		return nil, logrus.Fatal(ctx, err, nil)
	}

	// page infra. dependencies
	pageRepo := pagerepository.NewMongo(cfg, logrus, mongodb)
	pageParser := pageparser.NewHTML(logrus)
	pageScanner := pagescanner.NewHTTP(cfg, pageParser, logrus)

	// page domain dependencies
	pageComparator := pagecomparator.NewEqual()
	pageConsumer := pageconsumer.NewParallel(logrus, pageRepo, pageComparator)
	pageProvider := pageprovider.NewChan(cfg, pageScanner)

	// task dependencies
	taskParser := taskparser.NewCSV(cfg, logrus)
	taskProvider, err := taskprovider.NewParallel(ctx, logrus, cfg, taskParser)
	if err != nil {
		return nil, logrus.Fatal(ctx, err, nil)
	}
	taskRunner := taskrunner.New(cfg, pageProvider, pageConsumer)
	taskConsumer := taskconsumer.NewParallel(logrus, cfg, taskRunner)

	// job dependencies
	jobRunner := jobrunner.New(logrus, cfg, taskConsumer, taskProvider)
	jobScheduler := jobscheduler.NewTicker(cfg)

	return &App{
		ctx:          ctx,
		cancel:       cancel,
		config:       cfg,
		logger:       logrus,
		loggerCancel: logrusCancel,
		jobRunner:    jobRunner,
		jobScheduler: jobScheduler,
		liveness:     liveness,
		mongo:        mongodb,
	}, nil
}

func (a *App) Start() {
	defer a.Stop()
	a.logger.InfoMsg(a.ctx, "starting the spider application", nil)

	cancelLiveness := a.liveness.Watch(a)
	defer cancelLiveness()

	for range a.jobScheduler.Manage(a.ctx) {
		a.jobRunner.Run(a.ctx)
	}
}

func (a *App) Stop() {
	a.logger.InfoMsg(a.ctx, "closing the spider application", nil)

	a.cancel()
	a.loggerCancel()
}

func (a *App) IsAlive() bool {
	if err := a.mongo.Client().Ping(a.ctx, nil); err != nil {
		return false
	}
	return true
}
