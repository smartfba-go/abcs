package snowflake

import (
	"errors"
	"time"
)

const (
	UNIX_EPOCH_OFFSET_MILLIS uint64 = 1577858400000 // First millisecond of the year 2020

	SEQ_BITS   uint32 = 10
	SHARD_BITS uint32 = 12
	TS_BITS    uint32 = 41

	SEQ_MASK   uint64 = (uint64(1) << SEQ_BITS) - 1
	SHARD_MASK uint64 = (uint64(1) << SHARD_BITS) - 1
	TS_MASK    uint64 = (uint64(1) << TS_BITS) - 1

	SEQ_OFFSET   uint32 = 0
	SHARD_OFFSET uint32 = SEQ_OFFSET + SEQ_BITS
	TS_OFFSET    uint32 = SHARD_OFFSET + SHARD_BITS
)

type Pool struct {
	shard uint64
	seq   uint64
	ts    uint64
}

func NewPool(shard uint64) (*Pool, error) {
	if (shard & SHARD_MASK) != shard {
		return nil, ErrInvalidShard
	}

	return &Pool{
		shard: shard,
		seq:   0,
		ts:    0,
	}, nil
}

func (pool *Pool) IsAvailable() bool {
	return ((pool.seq + 1) & SEQ_MASK) != 0
}

func (pool *Pool) NewID() (ID, error) {
	now := time.Now()

	if err := pool.SetTime(now); err != nil {
		return 0, err
	}

	return pool.UnsafeNextID()
}

func (pool *Pool) UnsafeNextID() (ID, error) {
	if pool.seq > SEQ_MASK {
		availableAtMillis := UNIX_EPOCH_OFFSET_MILLIS + uint64(pool.ts>>uint64(TS_OFFSET)) + 1

		return 0, &PoolUnavailableError{
			AvailbleAt: time.UnixMilli(int64(availableAtMillis)),
		}
	}

	v := (pool.ts << TS_OFFSET) | (pool.shard << SHARD_OFFSET) | (pool.seq << SEQ_OFFSET)

	pool.seq += 1

	return ID(v), nil
}

func (pool *Pool) SetTime(t time.Time) error {
	rawTs := uint64(t.UnixMilli()) - UNIX_EPOCH_OFFSET_MILLIS
	ts := rawTs & TS_MASK

	if rawTs != ts {
		panic("permanently unavailable: out of time")
	}

	if pool.ts > ts {
		availableAtMillis := UNIX_EPOCH_OFFSET_MILLIS + uint64(pool.ts>>uint64(TS_OFFSET)) + 1

		return &PoolUnavailableError{
			AvailbleAt: time.UnixMilli(int64(availableAtMillis)),
		}
	}

	if pool.ts != ts {
		pool.ts = ts
		pool.seq = 0
	}

	return nil
}

type PoolUnavailableError struct {
	AvailbleAt time.Time
}

func (e *PoolUnavailableError) Error() string { return "pool unavailable" }

func (e *PoolUnavailableError) Unwrap() error { return nil }

var (
	ErrInvalidShard = errors.New("invalid shard")
)
