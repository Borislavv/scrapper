package pagescrapper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type PageScrapper struct {
	config spiderconfiginterface.Config
}

func New(config spiderconfiginterface.Config) *PageScrapper {
	return &PageScrapper{
		config: config,
	}
}

func (s *PageScrapper) Scrape(url *url.URL) (*entity.Page, error) {
	userAgent := "Chrome/118.0.5993.70 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"

	data := map[string]string{
		"url":       url.String(),
		"userAgent": userAgent,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://host.docker.internal:3000/scrape", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response from server:", string(body))

	return nil, errors.New("test error")
}
