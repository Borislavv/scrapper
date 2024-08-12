package spider

import (
	"context"
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
	taskprovider "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/task/provider"
	taskrunner "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/task/runner"
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

	// page dependencies
	pageRepo := pagerepository.New()
	pageSaver := pagesaver.New(pageRepo)
	pageFinder := pagefinder.New(pageRepo)
	pageScrapper := pagescrapper.New(ctx, cfg)
	pageComparator := pagecomparator.New()
	// task dependencies
	taskProvider := taskprovider.New(cfg)
	taskRunner := taskrunner.New(pageSaver, pageFinder, pageScrapper, pageComparator)
	// job dependencies
	jobRunner := jobrunner.New(cfg, taskRunner, taskProvider)
	jobScheduler := jobscheduler.New(cfg)

	return &Spider{
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
