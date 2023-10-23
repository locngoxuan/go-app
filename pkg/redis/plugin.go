package redis

import (
	"fmt"
	"goapp/helper"
	"goapp/pkg/interfaces"
	"goapp/plugin"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func NewRedis(ctx *cli.Context) interfaces.Redis {
	opt := RedisOption{
		Addr:     ctx.String(interfaces.FlagRedisAddr),
		Password: ctx.String(interfaces.FlagRedisPassword),
		Database: ctx.Int(interfaces.FlagRedisDatabase),
	}
	return &RedisPlugin{
		RedisOption: opt,
	}
}

type RedisOption struct {
	Addr     string
	Password string
	Database int
}

type RedisPlugin struct {
	plugin.EmptyPlugin
	RedisOption
	client *redis.Client
}

// Client implements interfaces.Redis.
func (r *RedisPlugin) Client() *redis.Client {
	return r.client
}

func (r *RedisPlugin) Start() error {
	log.Info().Msg("start redis plugin...")
	if helper.IsBlank(r.Addr) {
		return fmt.Errorf("address must not empty")
	}
	if r.Database < 0 || r.Database > 15 {
		return fmt.Errorf("database index must be a value between 0 and 15")
	}
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password, // no password set
		DB:       r.Database, // use default DB
	})
	return nil
}

func (r *RedisPlugin) Stop() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}
