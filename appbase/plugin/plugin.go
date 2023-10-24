package plugin

import (
	"reflect"

	"github.com/rs/zerolog/log"
)

type Plugin interface {
	Start() error
	Stop() error
}

type EmptyPlugin struct {
}

func (e *EmptyPlugin) Start() error {
	return nil
}
func (e *EmptyPlugin) Stop() error {
	return nil
}

func Start(plugins []interface{}) error {
	for _, p := range plugins {
		calls := reflect.ValueOf(p).Call([]reflect.Value{})
		err := calls[0].Interface().(Plugin).Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func Stop(plugins []interface{}) {
	for _, p := range plugins {
		calls := reflect.ValueOf(p).Call([]reflect.Value{})
		err := calls[0].Interface().(Plugin).Stop()
		if err != nil {
			log.Error().Err(err).Msg("failed to stop plugin")
		}
	}
}
