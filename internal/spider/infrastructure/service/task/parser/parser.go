package taskparser

import (
	"encoding/csv"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/util"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	"log"
	"net/url"
	"os"
)

type TaskParser struct {
	config spiderconfiginterface.Config
}

func New(config spiderconfiginterface.Config) *TaskParser {
	return &TaskParser{config: config}
}

func (p *TaskParser) ParseURLs() ([]*url.URL, error) {
	filepath, err := util.Path(p.config.GetURLsFilepath())
	if err != nil {
		log.Println("TaskParser.ParseURLs(): " + err.Error())
		return nil, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Println("TaskParser.ParseURLs(): " + err.Error())
		return nil, err
	}
	defer func() { _ = file.Close() }()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Println("TaskParser.ParseURLs(): " + err.Error())
		return nil, err
	}

	rawURLs := make([]string, len(records))
	for _, record := range records {
		if len(record) > 0 {
			rawURLs = append(rawURLs, record[0])
		}
	}

	URLs := make([]*url.URL, len(rawURLs))
	for _, rawURL := range rawURLs {
		URL, err := url.Parse(rawURL)
		if err != nil {
			log.Println("TaskParser.ParseURLs(): " + err.Error())
			return nil, err
		}
		URLs = append(URLs, URL)
	}

	return URLs, nil
}
