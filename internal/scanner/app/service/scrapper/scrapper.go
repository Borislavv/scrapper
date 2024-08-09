package scrapper

import (
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"time"
)

const (
	// Адрес Selenium сервера
	seleniumURL = "http://host.docker.internal:4444/wd/hub"
)

type Scrapper struct {
}

func New() *Scrapper {
	return &Scrapper{}
}

func (s *Scrapper) Scrape() {
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
		log.Fatalf("error connecting to the WebDriver: %v", err)
	}
	defer wd.Quit()

	if err := wd.SetImplicitWaitTimeout(5 * time.Second); err != nil {
		log.Fatalf("failed to set implicit wait timeout: %v", err)
	}

	if err := wd.Get("https://melbet-ar1.com/es/live/esports"); err != nil {
		log.Fatalf("failed to load page: %s\n", err)
	}

	// title
	title, err := wd.Title()
	if err != nil {
		log.Fatalf("failed to get page title: %s\n", err)
	}
	fmt.Printf("page title: %s\n", title)

	// description
	description, err := wd.FindElement(selenium.ByXPATH, `//meta[@name="description"]`)
	if err != nil {
		log.Printf("failed to find description: %s\n", err)
	} else {
		desc, err := description.GetAttribute("content")
		if err != nil {
			log.Printf("failed to get description content: %s\n", err)
		} else {
			fmt.Printf("page description: %s\n", desc)
		}
	}

	// H1
	h1Element, err := wd.FindElement(selenium.ByTagName, "h1")
	if err != nil {
		log.Printf("failed to find H1 tag: %s\n", err)
	} else {
		h1Text, err := h1Element.Text()
		if err != nil {
			log.Printf("failed to get H1 text: %s\n", err)
		} else {
			fmt.Printf("H1 tag text: %s\n", h1Text)
		}
	}

	// html
	_, err = wd.PageSource()
	if err != nil {
		log.Fatalf("failed to get page source: %s\n", err)
	}
	fmt.Println("page HTML: received")
	//fmt.Println(html)

	// logs
	consoleLogs, err := wd.Log("browser")
	if err != nil {
		log.Fatalf("failed to get console logs: %s\n", err)
	}
	for _, consoleLog := range consoleLogs {
		fmt.Printf("console log: %s\n", consoleLog)
	}
}
