package main

import (
	"github.com/glebarez/sqlite"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func Dialector(ctx *cli.Context) gorm.Dialector {
	return sqlite.Open("goapp.db")
}
func GormOpts() []gorm.Option {
	return []gorm.Option{
		&gorm.Config{
			SkipDefaultTransaction: true,
		},
	}
}
