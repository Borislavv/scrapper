package pagescrapper

import (
	"encoding/json"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	"github.com/tebeka/selenium"
	"log"
	"net/url"
	"runtime"
	"sync"
)

const seleniumURL = "http://host.docker.internal:4444/wd/hub"

type PageScrapper struct {
	config spiderconfiginterface.Config
	wdPool *sync.Pool
}

func New(config spiderconfiginterface.Config) *PageScrapper {
	return &PageScrapper{
		config: config,
		wdPool: &sync.Pool{
			New: func() interface{} {
				caps := selenium.Capabilities{
					"browserName": "chrome",
					"goog:chromeOptions": map[string]interface{}{
						"args": []string{
							"--no-sandbox",
							"--disable-dev-shm-usage",
							"--headless",
							"--enable-logging",
							"--v=1",
							"--enable-precise-memory-info",
							"--disable-popup-blocking",
							"--disable-default-apps",
							"--remote-debugging-port=9222",
						},
					},
					"goog:loggingPrefs": map[string]interface{}{
						"browser":     "ALL",
						"performance": "ALL",
					},
				}

				wd, err := selenium.NewRemote(caps, seleniumURL)
				if err != nil {
					log.Println("PageScrapper: " + err.Error())
					return nil
				}

				runtime.SetFinalizer(wd, func(w selenium.WebDriver) {
					if err = w.Quit(); err != nil {
						log.Println("PageScrapper: " + err.Error())
					}
				})

				return wd
			},
		},
	}
}

func (s *PageScrapper) Scrape(url *url.URL) (*entity.Page, error) {
	wd, ok := s.wdPool.Get().(selenium.WebDriver)
	if !ok {
		err := errors.New("cast to selenium.WebDriver failed (probably selenium.WebDriver creation failed)")
		log.Println("PageScrapper: " + err.Error())
		return nil, err
	}

	if err := wd.SetImplicitWaitTimeout(s.config.GetTimeoutPerURL()); err != nil {
		log.Printf("PageScrapper: failed to set implicit wait timeout: %v\n", err)
		return nil, err
	}

	if err := wd.Get(url.String()); err != nil {
		log.Printf("PageScrapper: failed to load page: %s\n", err)
		return nil, err
	}

	page := &entity.Page{
		URL: url.String(),
	}

	// title
	title, err := wd.Title()
	if err != nil {
		log.Printf("PageScrapper: failed to get page title: %s\n", err)
		return nil, err
	}
	page.Title = title

	// description
	descriptionElement, err := wd.FindElement(selenium.ByXPATH, "//meta[@name='description']")
	if err != nil {
		log.Printf("PageScrapper: failed to find description: %s\n", err)
		return nil, err
	} else {
		description, err := descriptionElement.GetAttribute("content")
		if err != nil {
			log.Printf("PageScrapper: failed to get description content: %s\n", err)
			return nil, err
		} else {
			page.Description = description
		}
	}

	// H1
	h1Element, err := wd.FindElement(selenium.ByTagName, "h1")
	if err != nil {
		log.Printf("PageScrapper: failed to find H1 tag: %s\n", err)
		return nil, err
	} else {
		h1, err := h1Element.Text()
		if err != nil {
			log.Printf("PageScrapper: failed to get H1 text: %s\n", err)
			return nil, err
		} else {
			page.H1 = h1
		}
	}

	// html
	html, err := wd.PageSource()
	if err != nil {
		log.Printf("PageScrapper: failed to get page source: %s\n", err)
		return nil, err
	}
	page.HTML = html

	// Сканирование сети (Network)
	logs, err := wd.Log("performance")
	if err != nil {
		log.Printf("PageScrapper: failed to get performance logs: %s\n", err)
		return nil, err
	}

	networkLogs := make(map[string]*entity.NetworkLog)

	for _, logEntry := range logs {
		var message struct {
			Message struct {
				Method string                 `json:"method"`
				Params map[string]interface{} `json:"params"`
			} `json:"message"`
		}
		if err = json.Unmarshal([]byte(logEntry.Message), &message); err != nil {
			log.Printf("PageScrapper: failed to unmarshal log entry: %s\n", err)
			continue
		}

		requestID, hasRequestID := message.Message.Params["requestId"].(string)
		if !hasRequestID {
			continue
		}

		switch message.Message.Method {
		case "Network.requestWillBeSent":
			request := message.Message.Params["request"].(map[string]interface{})
			reqUrl := request["url"].(string)
			headers := request["headers"].(map[string]interface{})

			if _, exists := networkLogs[requestID]; !exists {
				networkLogs[requestID] = &entity.NetworkLog{}
			}
			networkLogs[requestID].URL = reqUrl
			networkLogs[requestID].RequestHeaders = headers

		case "Network.responseReceived":
			if _, exists := networkLogs[requestID]; !exists {
				networkLogs[requestID] = &entity.NetworkLog{}
			}
			response := message.Message.Params["response"].(map[string]interface{})
			headers := response["headers"].(map[string]interface{})
			status := int(response["status"].(float64))

			networkLogs[requestID].ResponseHeaders = headers
			networkLogs[requestID].StatusCode = status
		}
	}

	// console logs
	consoleLogs, err := wd.Log("browser")
	if err != nil {
		log.Printf("PageScrapper: failed to get console logs: %s\n", err)
		return nil, err
	}
	logMessages := make([]string, 0, len(consoleLogs))
	for _, consoleLog := range consoleLogs {
		logMessages = append(logMessages, consoleLog.Message)
	}

	return page, nil
}
