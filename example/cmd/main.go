package main

import (
	"go-app/appbase/app"
	"go-app/appbase/di"
	"go-app/appbase/pkg/http"
	"go-app/appbase/pkg/interfaces"
	"go-app/appbase/pkg/repository"
	"go-app/example/entity/model"
	"go-app/example/features/orders"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	version = "devel"
	binFile = "go-app"
)

func main() {
	opt := app.GoAppOpt{
		Name:      binFile,
		Copyright: "Loc Ngo <xuanloc0511@gmail.com>",
		Version:   version,
		Args:      os.Args,
		Flags:     appFlags(),
		PreStart: func(ctx *cli.Context) error {
			di.GetRepository().SetDialector(Dialector(ctx))
			di.GetRepository().SetOptions(GormOpts())
			di.GetHttp().SetRouter(Apis())
			return nil
		},
		Ready: func(ctx *cli.Context) error {
			di.GetRepository().Database().AutoMigrate(&model.Order{})
			return nil
		},
		Injection: []interface{}{
			orders.NewOrderApi,
			orders.NewOrderService,
			orders.NewOrderRepository,
			repository.NewRepository,
			http.NewHttp,
		},
		Plugins: []interface{}{
			di.GetRepository,
			di.GetHttp,
		},
	}
	app, err := app.NewApp(opt)
	if err != nil {
		log.Error().Err(err).Msg("failed to create application")
	}
	err = app.Run()
	if err != nil {
		log.Error().Err(err).Msg("failed to run application")
	}
}

func appFlags() []cli.Flag {
	flags := make([]cli.Flag, 0)
	flags = append(flags, interfaces.HttpCliArg()...)
	flags = append(flags, interfaces.RepositoryCliArg()...)
	return flags
}
