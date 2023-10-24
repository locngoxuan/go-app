package app

import (
	"fmt"
	"go-app/appbase/di"
	"go-app/appbase/plugin"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type GoApp struct {
	opt GoAppOpt
	app *cli.App
}

type Configure func(ctx *cli.Context) error

type GoAppOpt struct {
	Version      string
	Copyright    string
	Name         string
	Injection    []interface{}
	Plugins      []interface{}
	Args         []string
	Flags        []cli.Flag
	PreStart     Configure
	Ready        Configure
	ConfigureLog Configure
}

func (a GoApp) Run() error {
	a.app.Action = a.run
	return a.app.Run(a.opt.Args)
}

func (a GoApp) setupInjection(ctx *cli.Context) error {
	// initialize dependency container and injection
	constructors := []interface{}{
		func() *cli.Context {
			return ctx
		},
	}
	constructors = append(constructors, a.opt.Injection...)

	return di.Initialize(constructors...)
}

// run handles logic before starting application. It also defines trigger for closing application
func (a GoApp) run(ctx *cli.Context) error {
	log.Info().Interface("args", os.Args).Msg("execute binary file with args")
	a.opt.ConfigureLog(ctx)
	start := time.Now()
	log.Info().Msg("application is starting...")

	err := a.setupInjection(ctx)
	if err != nil {
		return err
	}
	//setup APIs
	err = a.opt.PreStart(ctx)
	if err != nil {
		return err
	}

	//start plugins
	err = plugin.Start(a.opt.Plugins)
	if err != nil {
		return err
	}
	//migration
	if a.opt.Ready != nil {
		err = a.opt.Ready(ctx)
		if err != nil {
			return err
		}
	}
	log.Info().Int64("duration", time.Now().Sub(start).Milliseconds()).Msg("application has been up successful")
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGINT,
	)
	<-signalChannel
	signal.Stop(signalChannel)
	close(signalChannel)
	err = a.afterShutdown()
	if err != nil {
		log.Warn().Err(err).Msg("get error while shuting down application")
	}
	log.Info().Msg("application has been shutdown")
	return nil
}

// afterShutdown handles specific logic while shutting down application
func (a GoApp) afterShutdown() error {
	log.Info().Msg("shutting down application...")
	plugin.Stop(a.opt.Plugins)
	return nil
}

func NewApp(opt GoAppOpt) (goApp *GoApp, err error) {
	if opt.Injection == nil {
		return nil, fmt.Errorf("")
	}
	if opt.Plugins == nil {
		return nil, fmt.Errorf("")
	}
	if opt.ConfigureLog == nil {
		opt.ConfigureLog = setupZeroLog
	}
	app := &cli.App{
		Name:      opt.Name,
		Version:   opt.Version,
		Copyright: opt.Copyright,
		Flags:     opt.Flags,
	}
	return &GoApp{
		opt: opt,
		app: app,
	}, nil
}

func setupZeroLog(ctx *cli.Context) error {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	return nil
}
