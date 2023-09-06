package temporalworker

import (
	"context"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"temporalworker",
	fx.Provide(New),
)

type Params struct {
	fx.In

	Config Config
	Client client.Client
}

type Result struct {
	fx.Out

	Worker worker.Worker
}

type Config struct {
	TaskQueue string
	Options   worker.Options
}

func New(lc fx.Lifecycle, p Params) (Result, error) {
	w := worker.New(p.Client, p.Config.TaskQueue, p.Config.Options)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return w.Start()
		},
		OnStop: func(ctx context.Context) error {
			w.Stop()
			return nil
		},
	})

	return Result{
		Worker: w,
	}, nil
}
