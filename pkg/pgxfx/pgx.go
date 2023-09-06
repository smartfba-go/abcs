package pgxfx

import (
	"context"

	pgxpoolv4 "github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Module("pgx", fx.Provide(NewPool))

type PoolParams struct {
	fx.In

	Context context.Context `optional:"true"`
	Config  *Config
}

type PoolResult struct {
	fx.Out

	Pool   *pgxpool.Pool
	PoolV4 *pgxpoolv4.Pool
}

type Config struct {
	// The Data Source Name (DSN) or URI
	DSN string

	User     string
	Password string
}

func NewPool(lc fx.Lifecycle, p PoolParams) (PoolResult, error) {
	config, err := pgxpool.ParseConfig(p.Config.DSN)
	if err != nil {
		return PoolResult{}, err
	}

	configV4, err := pgxpoolv4.ParseConfig(p.Config.DSN)
	if err != nil {
		return PoolResult{}, err
	}

	if len(p.Config.User) > 0 {
		config.ConnConfig.Config.User = p.Config.User
		configV4.ConnConfig.Config.User = p.Config.User
	}

	if len(p.Config.Password) > 0 {
		config.ConnConfig.Config.Password = p.Config.Password
		configV4.ConnConfig.Config.Password = p.Config.Password
	}

	ctx := p.Context
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancelF := context.WithCancel(ctx)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		cancelF()

		return PoolResult{}, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return pool.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			defer cancelF()

			pool.Close()

			return nil
		},
	})

	poolV4, err := pgxpoolv4.ConnectConfig(ctx, configV4)
	if err != nil {
		cancelF()

		return PoolResult{}, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return poolV4.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			defer cancelF()

			poolV4.Close()

			return nil
		},
	})

	return PoolResult{
		Pool:   pool,
		PoolV4: poolV4,
	}, nil
}
