package pagescrapper

import (
	"context"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	"github.com/tebeka/selenium"
	"log"
	"net/url"
)

const seleniumURL = "http://host.docker.internal:4444/wd/hub"

type PageScrapper struct {
	config spiderconfiginterface.Config
}

func New(config spiderconfiginterface.Config) *PageScrapper {
	return &PageScrapper{config: config}
}

func (s *PageScrapper) Scrape(ctx context.Context, url url.URL) (*entity.Page, error) {
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"goog:chromeOptions": map[string]interface{}{
			"args": []string{
				"--no-sandbox",
				"--disable-dev-shm-usage",
				"--headless",
			},
		},
		"goog:loggingPrefs": map[string]interface{}{
			"browser":     "ALL",
			"performance": "ALL",
		},
	}

	wd, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		log.Printf("error connecting to the WebDriver: %v\n", err)
		return nil, err
	}
	defer func() { _ = wd.Quit() }()

	if err = wd.SetImplicitWaitTimeout(s.config.GetTimeoutPerURL()); err != nil {
		log.Printf("failed to set implicit wait timeout: %v", err)
		return nil, err
	}

	if err = wd.Get(url.String()); err != nil {
		log.Printf("failed to load page: %s\n", err)
		return nil, err
	}

	page := &entity.Page{}

	// title
	title, err := wd.Title()
	if err != nil {
		log.Printf("failed to get page title: %s\n", err)
		return nil, err
	}
	page.Title = title

	// description
	descriptionElement, err := wd.FindElement(selenium.ByXPATH, `//meta[@name="description"]`)
	if err != nil {
		log.Printf("failed to find description: %s\n", err)
		return nil, err
	} else {
		description, err := descriptionElement.GetAttribute("content")
		if err != nil {
			log.Printf("failed to get description content: %s\n", err)
			return nil, err
		} else {
			page.Description = description
		}
	}

	// H1
	h1Element, err := wd.FindElement(selenium.ByTagName, "h1")
	if err != nil {
		log.Printf("failed to find H1 tag: %s\n", err)
		return nil, err
	} else {
		h1, err := h1Element.Text()
		if err != nil {
			log.Printf("failed to get H1 text: %s\n", err)
			return nil, err
		} else {
			page.H1 = h1
		}
	}

	// html
	html, err := wd.PageSource()
	if err != nil {
		log.Printf("failed to get page source: %s\n", err)
		return nil, err
	}
	page.HTML = html

	// logs
	consoleLogs, err := wd.Log("browser")
	if err != nil {
		log.Printf("failed to get console logs: %s\n", err)
		return nil, err
	}
	logs := make([]string, 0, len(consoleLogs))
	for _, consoleLog := range consoleLogs {
		logs = append(logs, consoleLog.Message)
	}

	return page, nil
}
