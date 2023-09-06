package snowflaketemporal

import (
	"context"

	"go.smartfba.io/abcs/snowflake"
	snowflakeservice "go.smartfba.io/abcs/snowflake/service"
)

type Activities struct {
	Service snowflakeservice.SnowflakeService
}

func (a *Activities) NewSnowflakeActivity(ctx context.Context) (snowflake.ID, error) {
	id, err := a.Service.NewID(ctx)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *Activities) BatchNewSnowflakeActivity(ctx context.Context, n int32) ([]snowflake.ID, error) {
	ids, err := a.Service.NewIDs(ctx, n)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
