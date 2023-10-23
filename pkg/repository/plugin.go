package repository

import (
	"fmt"
	"goapp/pkg/interfaces"
	"goapp/plugin"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewRepository(ctx *cli.Context) interfaces.Repository {
	opt := GormRepoOption{}
	return &GormRepository{
		GormRepoOption: opt,
	}
}

type GormRepoOption struct {
	DSN          string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifeTime  time.Duration
	MaxIdleTime  time.Duration
	MaxExecTime  time.Duration
	LogLevel     string
}

type GormRepository struct {
	plugin.EmptyPlugin
	db *gorm.DB
	GormRepoOption
	dialector gorm.Dialector
	opts      []gorm.Option
}

// SetDialector implements interfaces.Repository.
func (g *GormRepository) SetDialector(dialector gorm.Dialector) {
	g.dialector = dialector
}

// SetOptions implements interfaces.Repository.
func (g *GormRepository) SetOptions(opts []gorm.Option) {
	g.opts = opts
}

// Database implements interfaces.Repository.
func (g *GormRepository) Database() *gorm.DB {
	return g.db
}

func (g *GormRepository) Start() error {
	log.Info().Msg("start repository plugin...")
	logLevel := logger.Silent
	switch strings.ToLower(g.LogLevel) {
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info", "debug":
		logLevel = logger.Info
	default:
		logLevel = logger.Silent
	}
	opts := []gorm.Option{
		&gorm.Config{
			Logger: logger.New(&log.Logger, logger.Config{
				SlowThreshold:             g.MaxExecTime, // Slow SQL threshold
				LogLevel:                  logLevel,      // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,         // Disable color
			}),
		},
	}
	if g.opts != nil {
		opts = append(opts, g.opts...)
	}
	var err error
	g.db, err = gorm.Open(g.dialector, opts...)
	if err != nil {
		return fmt.Errorf("failed to initialize db session, err: %v", err)
	}
	sqlDB, err := g.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB instance, err: %v", err)
	}
	if g.MaxIdleConns < 0 || g.MaxOpenConns < 0 {
		return fmt.Errorf("max idle connection and max open connection must are positive number")
	}
	if g.MaxIdleConns > g.MaxOpenConns {
		return fmt.Errorf("number of open connection must larger than number of idle connection")
	}
	sqlDB.SetMaxIdleConns(g.MaxIdleConns)
	sqlDB.SetMaxOpenConns(g.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(g.MaxLifeTime)
	sqlDB.SetConnMaxIdleTime(g.MaxIdleTime)
	log.Info().Msg("open connection to database successful")
	return nil
}

func (g *GormRepository) Stop() error {
	if g.db != nil {
		sqlDB, err := g.db.DB()
		if err != nil {
			return err
		}
		sqlDB.Close()
	}
	return nil
}
