package taskparser

import (
	"encoding/csv"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/util"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	"log"
	"net/url"
	"os"
)

type CSV struct {
	config spiderconfiginterface.Config
}

func NewCSV(config spiderconfiginterface.Config) *CSV {
	return &CSV{config: config}
}

func (p *CSV) Parse() ([]*url.URL, error) {
	filepath, err := util.Path(p.config.GetURLsFilepath())
	if err != nil {
		log.Println("CSV.Parse(): " + err.Error())
		return nil, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Println("CSV.Parse(): " + err.Error())
		return nil, err
	}
	defer func() { _ = file.Close() }()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Println("CSV.Parse(): " + err.Error())
		return nil, err
	}

	rawURLs := make([]string, 0, len(records))
	for _, record := range records {
		if len(record) > 0 {
			rawURLs = append(rawURLs, record[0])
			continue
		}
	}

	URLs := make([]*url.URL, 0, len(rawURLs))
	for _, rawURL := range rawURLs {
		URL, err := url.Parse(rawURL)
		if err != nil {
			log.Println("CSV.Parse(): " + err.Error())
			return nil, err
		}
		URLs = append(URLs, URL)
	}

	return URLs, nil
}
