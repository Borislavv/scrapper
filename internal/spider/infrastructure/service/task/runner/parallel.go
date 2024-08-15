package taskrunner

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	pagerepositoryinterface "github.com/Borislavv/scrapper/internal/spider/domain/repository/interface"
	pagecomparatorinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/comparator/interface"
	pagefinderinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/finder/interface"
	pagesaverinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/saver/interface"
	pagescannerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/scanner/interface"
	"log"
	"net/url"
	"sync"
)

type TaskRunner struct {
	config     spiderconfiginterface.Config
	saver      pagesaverinterface.PageSaver
	finder     pagefinderinterface.PageFinder
	scanner    pagescannerinterface.PageScanner
	comparator pagecomparatorinterface.PageComparator
}

func New(
	config spiderconfiginterface.Config,
	saver pagesaverinterface.PageSaver,
	finder pagefinderinterface.PageFinder,
	scanner pagescannerinterface.PageScanner,
	comparator pagecomparatorinterface.PageComparator,
) *TaskRunner {
	return &TaskRunner{
		config:     config,
		saver:      saver,
		finder:     finder,
		scanner:    scanner,
		comparator: comparator,
	}
}

func (r *TaskRunner) Run(ctx context.Context, wg *sync.WaitGroup, url *url.URL) {
	//defer log.Println("TaskRunner: closed")
	//log.Println("TaskRunner: started")

	defer wg.Done()

	errsCh := make(chan error, 1)
	pagesCh := make(chan *entity.Page, 1)

	cwg := &sync.WaitGroup{}
	defer cwg.Wait()
	cwg.Add(1)
	go r.consumeErrors(cwg, errsCh)
	cwg.Add(1)
	go r.consumePages(ctx, cwg, url, pagesCh)

	defer func() { close(errsCh); close(pagesCh) }()

	pwg := &sync.WaitGroup{}
	defer pwg.Wait()
	pwg.Add(1)
	go r.providePagesAndErrors(ctx, pwg, url, pagesCh, errsCh)
}

func (r *TaskRunner) providePagesAndErrors(
	ctx context.Context,
	wg *sync.WaitGroup,
	url *url.URL,
	pagesCh chan *entity.Page,
	errsCh chan error,
) {
	defer wg.Done()

	for _, userAgent := range r.config.GetUserAgents() {
		wg.Add(1)
		go r.scanner.Scan(ctx, wg, url, userAgent, pagesCh, errsCh, r.config.GetRequestRetries())
	}
}

func (r *TaskRunner) consumeErrors(wg *sync.WaitGroup, errsCh <-chan error) {
	defer wg.Done()

	for err := range errsCh {
		log.Println("--------TaskRunner: " + err.Error())
	}
}

func (r *TaskRunner) consumePages(ctx context.Context, wg *sync.WaitGroup, url *url.URL, pagesCh <-chan *entity.Page) {
	defer wg.Done()

	for cur := range pagesCh {
		log.Println("=========TaskRunner: " + cur.URL)
		prev, err := r.finder.FindByURL(ctx, url)
		if err != nil {
			if errors.Is(err, pagerepositoryinterface.NotFoundError) {
				if err = r.saver.Save(ctx, cur); err != nil {
					log.Println("TaskRunner: " + err.Error())
				}
				log.Printf("TaskRunner: page with url %s was saved at first time.\n", url.String())
				return
			}

			log.Println("TaskRunner: " + err.Error())
			return
		}

		if !r.comparator.IsEquals(cur, prev) {
			log.Println("TaskRunner: blinking detected for url: " + cur.URL)
			if err = r.saver.Save(ctx, cur); err != nil {
				log.Println("TaskRunner: " + err.Error())
			}
		} else {
			log.Println("TaskRunner: pages are equal for url: " + cur.URL)
		}
	}
}
