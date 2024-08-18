package taskparser

import (
	"context"
	"encoding/csv"
	"gitlab.xbet.lan/web-backend/php/spider/internal/shared/infrastructure/util"
	spiderconfiginterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/app/config/interface"
	taskparserinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/task/parser/interface"
	logger "gitlab.xbet.lan/web-backend/php/spider/internal/spider/infrastructure/logger/interface"
	"net/url"
	"os"
)

type CSV struct {
	config spiderconfiginterface.Configurator
	logger logger.Logger
}

// NewCSV is a constructor of CSV parser.
func NewCSV(config spiderconfiginterface.Configurator, logger logger.Logger) *CSV {
	return &CSV{config: config, logger: logger}
}

func (p *CSV) Parse(ctx context.Context) ([]*url.URL, error) {
	filepath, err := util.Path(p.config.GetURLsFilepath())
	if err != nil {
		return nil, p.logger.Error(ctx, taskparserinterface.BuildFilepathError, logger.Fields{
			"path": p.config.GetURLsFilepath(),
			"err":  err.Error(),
		})
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, p.logger.Error(ctx, taskparserinterface.OpenFileError, logger.Fields{
			"path": p.config.GetURLsFilepath(),
			"err":  err.Error(),
		})
	}
	defer func() { _ = file.Close() }()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, p.logger.Error(ctx, taskparserinterface.ReadFileError, logger.Fields{
			"path": p.config.GetURLsFilepath(),
			"err":  err.Error(),
		})
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
			p.logger.WarningMsg(ctx, taskparserinterface.ParseURLError.Error(), logger.Fields{
				"path": p.config.GetURLsFilepath(),
				"err":  err.Error(),
				"url":  rawURL,
			})
			continue
		}
		URLs = append(URLs, URL)
	}

	return URLs, nil
}
