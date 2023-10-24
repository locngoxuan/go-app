package interfaces

import (
	"go-app/appbase/plugin"
	"net/http"

	"github.com/urfave/cli/v2"
)

type HTTP interface {
	plugin.Plugin
	SetRouter(router http.Handler)
}

const (
	FlagHttpBind = "http.bind"
	FlagHttpPort = "http.port"
)

func HttpCliArg() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    FlagHttpBind,
			EnvVars: []string{"HTTP_BIND"},
			Usage:   "",
			Value:   "0.0.0.0",
		},
		&cli.IntFlag{
			Name:    FlagHttpPort,
			EnvVars: []string{"HTTP_PORT"},
			Usage:   "",
			Value:   8080,
		},
	}
}
