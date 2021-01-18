package redis

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/rs/zerolog/log"
)

// Lock 获取锁
func Lock(ctx context.Context,k string, ttl time.Duration) (unlockfn func() error, err error) {

	lockermtx.Lock()
	defer lockermtx.Unlock()

	lk := &lock{
		ctx: ctx,
		ttl:      ttl,
		watchTTL: ttl - time.Millisecond*20,
	}
	if lk.ttl <= time.Millisecond*20 {
		lk.watchTTL = lk.ttl / 5 * 4
	}

RELOCK:
	lk.ctx, lk.cancelfn = context.WithCancel(context.Background())
	lk.options = &redislock.Options{

	}
	if lk.lock, err = locker.Obtain(lk.ctx,k, ttl, lk.options); nil != err {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, redislock.ErrNotObtained) {
			goto RELOCK
		}
		return
	}

	unlockfn = func() error {
		lk.cancelfn()
		return lk.lock.Release(lk.ctx)
	}
	go lk.watchdog()

	return
}

type lock struct {
	lock     *redislock.Lock
	ctx      context.Context
	cancelfn context.CancelFunc
	ttl      time.Duration
	options  *redislock.Options
	watchTTL time.Duration
}

func (lk *lock) watchdog() {

	timer := time.NewTimer(lk.watchTTL)

	for {
		select {
		case <-lk.ctx.Done():
			return
		case <-timer.C:
			if err := lk.lock.Refresh(lk.ctx,lk.ttl, lk.options); nil != err {
				log.Err(err).Str("key", lk.lock.Key()).Msgf("watchdog refresh")
				return
			}
			timer.Reset(lk.watchTTL)
		}
	}
}
