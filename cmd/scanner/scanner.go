package scanner

import (
	"context"
	"fmt"
	"github.com/Borislavv/scrapper/internal/scanner/app/service/scheduler"
	"github.com/Borislavv/scrapper/internal/scanner/app/service/scrapper"
	"sync"
)

type Scanner struct {
	scrapper  *scrapper.Scrapper
	scheduler *scheduler.Scheduler
}

func NewApp() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	s.scheduler = scheduler.New()
	s.scrapper = scrapper.New()

	s.Scan(ctx)
}

func (s *Scanner) Scan(ctx context.Context) {
	fmt.Println("start scanning...")
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.scheduler.Manage(ctx):
			fmt.Println("scrapping...")
			s.scrapper.Scrape()
		}
	}
}
