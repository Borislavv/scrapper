package taskrunner

import (
	"context"
	"errors"
	pagerepositoryinterface "github.com/Borislavv/scrapper/internal/spider/domain/repository/interface"
	pagecomparatorinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/comparator/interface"
	pagefinderinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/finder/interface"
	pagesaverinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/saver/interface"
	pagescrapperinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/scrapper/interface"
	"log"
	"net/url"
	"sync"
)

type TaskRunner struct {
	finder     pagefinderinterface.PageFinder
	saver      pagesaverinterface.PageSaver
	scrapper   pagescrapperinterface.PageScrapper
	comparator pagecomparatorinterface.PageComparator
}

func New(
	saver pagesaverinterface.PageSaver,
	finder pagefinderinterface.PageFinder,
	scrapper pagescrapperinterface.PageScrapper,
	comparator pagecomparatorinterface.PageComparator,
) *TaskRunner {
	return &TaskRunner{
		saver:      saver,
		finder:     finder,
		scrapper:   scrapper,
		comparator: comparator,
	}
}

func (r *TaskRunner) Run(ctx context.Context, wg *sync.WaitGroup, url url.URL) error {
	defer wg.Done()

	cur, err := r.scrapper.Scrape(ctx, url)
	if err != nil {
		log.Println("TaskRunner: " + err.Error())
		return err
	}

	prev, err := r.finder.FindByURL(url)
	if err != nil {
		if errors.Is(err, pagerepositoryinterface.NotFoundError) {
			if err = r.saver.Save(cur); err != nil {
				log.Println("TaskRunner: " + err.Error())
				return err
			}
			log.Println("TaskRunner: saving: page was saved at first time.")
			return nil
		}

		log.Println("TaskRunner: " + err.Error())
		return err
	}

	if !r.comparator.IsEquals(cur, prev) {
		if err = r.saver.Save(cur); err != nil {
			log.Println("TaskRunner: " + err.Error())
			return err
		}
	}

	return nil
}
