package temporalclient

import (
	"context"

	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"temporalclient",
	fx.Provide(New),
)

type Params struct {
	fx.In

	Config Config
}

type Result struct {
	fx.Out

	Client client.Client
}

type Config struct {
	HostPort  string
	Namespace string
	Identity  string
}

func New(lc fx.Lifecycle, p Params) (Result, error) {
	cli, err := client.Dial(client.Options{
		HostPort:  p.Config.HostPort,
		Namespace: p.Config.Namespace,
		Identity:  p.Config.Namespace,
	})
	if err != nil {
		return Result{}, nil
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := cli.CheckHealth(ctx, &client.CheckHealthRequest{})

			return err
		},
		OnStop: func(ctx context.Context) error {
			cli.Close()

			return nil
		},
	})

	return Result{
		Client: cli,
	}, nil
}
