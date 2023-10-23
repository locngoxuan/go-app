package main

import (
	"goapp/pkg"
	"goapp/pkg/http"
	"goapp/pkg/interfaces"
	"goapp/pkg/repository"
	"goapp/plugin"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	version = "devel"
	binFile = "go-app"
)

func main() {
	app := &cli.App{
		Name:      binFile,
		Copyright: "Loc Ngo <xuanloc0511@gmail.com>",
		Version:   version,
		Flags:     appFlags(),
		Action:    run,
	}
	log.Info().Interface("args", os.Args).Msg("execute binary file with args")
	err := app.Run(os.Args)
	if err != nil {
		log.Error().Err(err).Msg("failed to run application")
	}
}

func appFlags() []cli.Flag {
	flags := make([]cli.Flag, 0)
	flags = append(flags, interfaces.HttpCliArg()...)
	return flags
}

var plugins = []interface{}{
	pkg.GetRepository,
	pkg.GetRedis,
	pkg.GetHttp,
}

// run handles logic before starting application. It also defines trigger for closing application
func run(ctx *cli.Context) error {
	start := time.Now()
	log.Info().Msg("application is starting...")
	//initialize dependency container and injection
	err := pkg.Initialize(pkg.Registration{
		Constructor: func() *cli.Context {
			return ctx
		},
	}, pkg.Registration{
		Constructor: http.NewHttp,
	}, pkg.Registration{
		Constructor: repository.NewRepository,
	})
	if err != nil {
		return err
	}
	//setup APIs
	pkg.GetRepository().SetDialector(Dialector(ctx))
	pkg.GetRepository().SetOptions(GormOpts())
	pkg.GetHttp().SetRouter(Apis())

	//start plugins
	err = plugin.Start(plugins)
	if err != nil {
		return err
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
	err = afterShutdown()
	if err != nil {
		log.Warn().Err(err).Msg("get error while shuting down application")
	}
	log.Info().Msg("application has been shutdown")
	return nil
}

// afterShutdown handles specific logic while shutting down application
func afterShutdown() error {
	log.Info().Msg("shutting down application...")
	plugin.Stop(plugins)
	return nil
}
