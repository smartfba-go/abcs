package snowflakesync

import (
	"context"
	"sync"
	"time"

	"go.smartfba.io/abcs/snowflake"
)

type Pool struct {
	p *snowflake.Pool

	mu sync.RWMutex
}

func NewPool(shard uint64) (*Pool, error) {
	p, err := snowflake.NewPool(shard)
	if err != nil {
		return nil, err
	}

	return &Pool{
		p: p,
	}, nil
}

func (pool *Pool) IsAvailable() bool {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.p.IsAvailable()
}

func (pool *Pool) NewID() (snowflake.ID, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.p.NewID()
}

func (pool *Pool) UnsafeNextID() (snowflake.ID, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.p.UnsafeNextID()
}

func (pool *Pool) SetTime(t time.Time) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.p.SetTime(t)
}

func (pool *Pool) NewIDs(ctx context.Context, n int) ([]snowflake.ID, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	ids := make([]snowflake.ID, n)

	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			return ids[:i], ctx.Err()
		default:
			id, err := pool.p.NewID()
			if err != nil {
				return ids[:i], err
			}
			ids[i] = id
		}
	}

	return ids, nil
}

func (pool *Pool) UnsafeNextIDs(ctx context.Context, n int) ([]snowflake.ID, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	ids := make([]snowflake.ID, n)

	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			return ids[:i], ctx.Err()
		default:
			id, err := pool.p.UnsafeNextID()
			if err != nil {
				return ids[:i], err
			}
			ids[i] = id
		}
	}

	return ids, nil
}
