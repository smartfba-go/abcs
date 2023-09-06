package redisfx

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var Module = fx.Module("redisfx", fx.Provide(NewClient))

type ClientParams struct {
	fx.In

	Context context.Context `optional:"true"`
	Config  *Config         `optional:"true"`
}

type ClientResult struct {
	fx.Out

	Client *redis.Client
}

type Config struct {
	Addr string
}

func NewClient(lc fx.Lifecycle, p ClientParams) (ClientResult, error) {
	config := p.Config
	if config == nil {
		config = &Config{}
	}

	client := redis.NewClient(&redis.Options{
		Addr: config.Addr,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx).Err()
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return ClientResult{
		Client: client,
	}, nil
}
