package interfaces

import (
	"go-app/appbase/plugin"
	"time"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

type Repository interface {
	plugin.Plugin
	Database() *gorm.DB
	SetDialector(dialector gorm.Dialector)
	SetOptions(opts []gorm.Option)
}

type DialectorCreator func() gorm.Dialector

const (
	FlagDbDSN          = "db.dsn"
	FlagDbMaxIdleConns = "db.con.maxIdle"
	FlagDbMaxOpenConns = "db.con.maxOpen"
	FlagDbMaxLifeTime  = "db.con.maxLifeTime"
	FlagDbMaxIdleTime  = "db.con.maxIdleTime"
	FlagDbMaxExecTime  = "db.con.maxExecTime"
	FlagDbLogLevel     = "db.logLevel"
)

func RepositoryCliArg() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    FlagDbDSN,
			EnvVars: []string{"DB_DSN"},
			Usage:   "",
			Value:   "",
		},
		&cli.IntFlag{
			Name:    FlagDbMaxIdleConns,
			EnvVars: []string{"DB_CON_MAXIDLE"},
			Usage:   "",
			Value:   1,
		},
		&cli.IntFlag{
			Name:    FlagDbMaxOpenConns,
			EnvVars: []string{"DB_CON_MAXOPEN"},
			Usage:   "",
			Value:   2,
		},
		&cli.DurationFlag{
			Name:    FlagDbMaxLifeTime,
			EnvVars: []string{"DB_CON_MAXLIFETIME"},
			Usage:   "",
			Value:   time.Duration(300 * time.Second),
		},
		&cli.DurationFlag{
			Name:    FlagDbMaxIdleTime,
			EnvVars: []string{"DB_CON_MAXIDLETIME"},
			Usage:   "",
			Value:   time.Duration(60 * time.Second),
		},
		&cli.DurationFlag{
			Name:    FlagDbMaxExecTime,
			EnvVars: []string{"DB_CON_MAXEXECTIME"},
			Usage:   "",
			Value:   time.Duration(30 * time.Second),
		},
		&cli.StringFlag{
			Name:    FlagDbLogLevel,
			EnvVars: []string{"DB_LOGLEVEL"},
			Usage:   "",
			Value:   "off",
		},
	}
}
