package pkg

import (
	"goapp/pkg/interfaces"

	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

var container = dig.New()

func Container() *dig.Container {
	return container
}

type Registration struct {
	Constructor interface{}
	Opts        []dig.ProvideOption
}

func Initialize(registrations ...Registration) error {
	var err error
	for _, registration := range registrations {
		if registration.Opts != nil {
			err = container.Provide(registration.Constructor, registration.Opts...)
		} else {
			err = container.Provide(registration.Constructor)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func GetHttp() interfaces.HTTP {
	var plugin interfaces.HTTP
	err := container.Invoke(func(element interfaces.HTTP) {
		plugin = element
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get http plugin")
	}
	return plugin
}

func GetRepository() interfaces.Repository {
	var plugin interfaces.Repository
	err := container.Invoke(func(element interfaces.Repository) {
		plugin = element
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get repository plugin")
	}
	return plugin
}

func GetRedis() interfaces.Redis {
	var plugin interfaces.Redis
	err := container.Invoke(func(element interfaces.Redis) {
		plugin = element
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get redis plugin")
	}
	return plugin
}
