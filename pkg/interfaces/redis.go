package interfaces

import (
	"goapp/plugin"

	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"
)

type Redis interface {
	plugin.Plugin
	Client() *redis.Client
}

const (
	FlagRedisAddr     = "redis.addr"
	FlagRedisPassword = "redis.password"
	FlagRedisDatabase = "redis.database"
)

func RedisCliArg() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    FlagRedisAddr,
			EnvVars: []string{"REDIS_ADDR"},
			Usage:   "",
			Value:   "localhost:6379",
		},
		&cli.StringFlag{
			Name:    FlagRedisPassword,
			EnvVars: []string{"REDIS_PASSWORD"},
			Usage:   "",
			Value:   "",
		},
		&cli.IntFlag{
			Name:    FlagRedisDatabase,
			EnvVars: []string{"REDIS_DATABASE"},
			Usage:   "",
			Value:   0,
		},
	}
}
