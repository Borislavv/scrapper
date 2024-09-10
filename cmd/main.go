package main

import (
	"context"
	"fmt"
	shared "github.com/Borislavv/scrapper/internal/shared/app"
	sharedconfig "github.com/Borislavv/scrapper/internal/shared/app/config"
	loggerinterface "github.com/Borislavv/scrapper/internal/shared/domain/service/logger/interface"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/liveness"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/logger"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/shutdown"
	spider "github.com/Borislavv/scrapper/internal/spider/app"
	"github.com/kelseyhightower/envconfig"
	"runtime"
)

func main() {
	cfg, err := sharedconfig.Load()
	if err != nil {
		panic(err)
	}

	output, cancelOutput, err := logger.NewOutput(cfg)
	if err != nil {
		panic(err)
	}
	defer cancelOutput()

	lgr, lgrCancel, err := logger.NewLogrus(cfg, output)
	if err != nil {
		panic(err)
	}
	defer lgrCancel()

	ctx, cancel := context.WithCancel(context.Background())

	gsh := shutdown.NewGraceful(ctx, cancel)

	lgr.InfoMsg(ctx, "starting the entire spider service", nil)

	if err = setMaxProc(ctx, lgr); err != nil {
		lgr.FatalMsg(ctx, "setting up GOMAXPROCS failed", loggerinterface.Fields{"err": err.Error()})
		return
	}

	livenessProbe := liveness.NewProbe(ctx)

	sh, err := shared.New(ctx, output, livenessProbe)
	if err != nil {
		lgr.FatalMsg(ctx, "starting the shared application failed", loggerinterface.Fields{"err": err.Error()})
		return
	}

	sp, err := spider.New(ctx, output, livenessProbe)
	if err != nil {
		lgr.FatalMsg(ctx, "starting the spider application failed", loggerinterface.Fields{"err": err.Error()})
		return
	}

	gsh.Add(1)
	go func() {
		defer gsh.Done()
		sp.Start()
	}()

	gsh.Add(1)
	go func() {
		defer gsh.Done()
		sh.Start()
	}()

	gsh.ListenCancelAndAwait()

	lgr.InfoMsg(ctx, "the entire spider service has been successfully shut down, exiting", nil)
}

func setMaxProc(ctx context.Context, lgr loggerinterface.Logger) error {
	type ProcessConfig struct {
		// GOMAXPROCS is a number of available system threads. If value is zero, then will be applied num. of CPU.
		GOMAXPROCS int `envconfig:"GOMAXPROCS" default:"0"`
	}

	cfg := &ProcessConfig{}
	if err := envconfig.Process("", cfg); err != nil {
		return err
	}

	if cfg.GOMAXPROCS == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(cfg.GOMAXPROCS)
	}

	lgr.DebugMsg(ctx, fmt.Sprintf("GOMAXPROCS: %d", runtime.GOMAXPROCS(0)), nil)

	return nil
}
