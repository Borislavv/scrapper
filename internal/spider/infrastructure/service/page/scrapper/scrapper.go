package pagescrapper

import (
	"encoding/json"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"net/url"
	"sync"
	"time"
)

// const seleniumURL = "http://selenium:4444/wd/hub"
const seleniumURL = "http://host.docker.internal:4444/wd/hub"

type PageScrapper struct {
	config spiderconfiginterface.Config
	WDs    *sync.Pool
}

func New(config spiderconfiginterface.Config) *PageScrapper {
	return &PageScrapper{
		config: config,
		WDs: &sync.Pool{
			New: func() any {
				chromeOptions := chrome.Capabilities{
					Args: []string{
						"--no-sandbox",
						"--disable-dev-shm-usage",
						"--enable-logging",
						"--v=1",
						"--enable-precise-memory-info",
						"--disable-popup-blocking",
						"--disable-default-apps",
						"--remote-debugging-port=9222",
						"--user-agent=Chrome/118.0.5993.70 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
					},
					W3C: true,
				}

				caps := selenium.Capabilities{
					"browserName":        "chrome",
					"goog:chromeOptions": chromeOptions,
					"goog:loggingPrefs": map[string]interface{}{
						"browser":     "ALL",
						"performance": "ALL",
					},
				}

				wd, err := selenium.NewRemote(caps, seleniumURL)
				if err != nil {
					log.Println("PageScrapper: " + err.Error())
				}

				return wd
			},
		},
	}
}

func (s *PageScrapper) getWD() (selenium.WebDriver, error) {
	chromeOptions := chrome.Capabilities{
		Args: []string{
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--enable-logging",
			"--v=1",
			"--enable-precise-memory-info",
			"--disable-popup-blocking",
			"--disable-default-apps",
			"--remote-debugging-port=9222",
			"--user-agent=Chrome/118.0.5993.70 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		},
		W3C: true,
	}

	caps := selenium.Capabilities{
		"browserName":        "chrome",
		"goog:chromeOptions": chromeOptions,
		"goog:loggingPrefs": map[string]interface{}{
			"browser":     "ALL",
			"performance": "ALL",
		},
	}

	wd, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		log.Println("PageScrapper: " + err.Error())
		return nil, err
	}

	return wd, nil
}

func (s *PageScrapper) getRecursiveWD(i int) (selenium.WebDriver, error) {
	wd, err := s.getWD()
	if err != nil {
		if i == 0 {
			return nil, errors.New("times exceeded")
		}
		time.Sleep(time.Second)
		return s.getRecursiveWD(i - 1)
	}
	return wd, nil
}

func (s *PageScrapper) Scrape(url *url.URL) (*entity.Page, error) {
	//wd := s.WDs.Get().(selenium.WebDriver)
	//defer s.WDs.Put(wd)
	wd, err := s.getRecursiveWD(10)
	if err != nil {
		log.Printf("PageScrapper: failed to init WebDriver: %s\n", err)
		return nil, err
	}
	defer wd.Close()

	//if err := wd.SetImplicitWaitTimeout(s.config.GetTimeoutPerURL()); err != nil {
	//	log.Printf("PageScrapper: failed to set implicit wait timeout: %v\n", err)
	//	return nil, err
	//}

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
	}
	page.Title = title

	// description
	descriptionElement, err := wd.FindElement(selenium.ByXPATH, "//meta[@name='description']")
	if err != nil {
		log.Printf("PageScrapper: failed to find description: %s\n", err)
	} else {
		description, err := descriptionElement.GetAttribute("content")
		if err != nil {
			log.Printf("PageScrapper: failed to get description content: %s\n", err)
		} else {
			page.Description = description
		}
	}

	// H1
	h1Element, err := wd.FindElement(selenium.ByTagName, "h1")
	if err != nil {
		log.Printf("PageScrapper: failed to find H1 tag: %s\n", err)
	} else {
		h1, err := h1Element.Text()
		if err != nil {
			log.Printf("PageScrapper: failed to get H1 text: %s\n", err)
		} else {
			page.H1 = h1
		}
	}

	// html
	html, err := wd.PageSource()
	if err != nil {
		log.Printf("PageScrapper: failed to get page source: %s\n", err)
	}
	page.HTML = html

	// Сканирование сети (Network)
	logs, err := wd.Log("performance")
	if err != nil {
		log.Printf("PageScrapper: failed to get performance logs: %s\n", err)
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
	}
	logMessages := make([]string, 0, len(consoleLogs))
	for _, consoleLog := range consoleLogs {
		logMessages = append(logMessages, consoleLog.Message)
	}

	return page, nil
}
