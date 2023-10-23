package main

import (
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func Dialector(ctx *cli.Context) gorm.Dialector {
	return nil
}
func GormOpts() []gorm.Option {
	return []gorm.Option{
		&gorm.Config{
			SkipDefaultTransaction: true,
		},
	}
}
