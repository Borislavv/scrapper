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
				"--disable-gpu",
				"--remote-debugging-port=9222",
			},
		},
	}
	wd, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		log.Fatalf("Error connecting to the WebDriver: %v", err)
	}
	defer wd.Quit()

	// Открыть URL
	if err := wd.Get("https://melbet-ar1.com/es/live/esports"); err != nil {
		log.Fatalf("Failed to load page: %s\n", err)
	}

	// Подождать, чтобы страница загрузилась
	time.Sleep(5 * time.Second)

	// Получить заголовок страницы
	title, err := wd.Title()
	if err != nil {
		log.Fatalf("Failed to get page title: %s\n", err)
	}
	fmt.Printf("Page title: %s\n", title)
}
