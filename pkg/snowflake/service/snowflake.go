package snowflakeservice

import (
	"context"
	"errors"

	"go.smartfba.io/abcs/pkg/failures"
	"go.smartfba.io/abcs/pkg/snowflake"
	snowflakesync "go.smartfba.io/abcs/pkg/snowflake/sync"
	"go.uber.org/fx"
)

var Module = fx.Module("snowflakeservice", fx.Provide(New))

func New(pool *snowflakesync.Pool) (SnowflakeService, error) {
	return &snowflakeService{
		Pool: pool,
	}, nil
}

type SnowflakeService interface {
	NewID(ctx context.Context) (snowflake.ID, error)
	NewIDs(ctx context.Context, n int32) ([]snowflake.ID, error)
}

type snowflakeService struct {
	Pool *snowflakesync.Pool
}

func (s *snowflakeService) NewID(ctx context.Context) (snowflake.ID, error) {
	id, err := s.Pool.NewID()
	if err != nil {
		return 0, convertError(err)
	}

	return id, nil
}

func (s *snowflakeService) NewIDs(ctx context.Context, n int32) ([]snowflake.ID, error) {
	ids, err := s.Pool.NewIDs(ctx, int(n))
	if err != nil {
		return nil, convertError(err)
	}

	return ids, nil
}

func convertError(err error) error {
	if err == nil {
		return nil
	}

	var (
		poolUnavailableError *snowflake.PoolUnavailableError
	)
	if errors.As(err, &poolUnavailableError) {
		return failures.Wrap(failures.Unavailable, "pool unavailable", err)
	}

	return err
}
