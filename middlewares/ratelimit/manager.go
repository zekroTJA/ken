package ratelimit

import (
	"fmt"
	"sync"
	"time"

	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/ken"
)

// Manager provides limiter instances by command instance,
// userID and guildID.
type Manager interface {
	// GetLimiter returns a limiter instance from the given
	// cmd instance, userID and guildID.
	//
	// cmd is guaranteed to also implement the LimitedCommand
	// interface when calling this function.
	GetLimiter(cmd ken.Command, userID, guildID string) *Limiter
}

type internalManager struct {
	limiters *timedmap.TimedMap
	pool     *sync.Pool
}

func newInternalManager() *internalManager {
	return &internalManager{
		limiters: timedmap.New(10 * time.Minute),
		pool: &sync.Pool{
			New: func() interface{} { return new(Limiter) },
		},
	}
}

func (c *internalManager) GetLimiter(cmd ken.Command, userID, guildID string) *Limiter {
	key := fmt.Sprintf("%s:%s:%s", cmd.Name(), guildID, userID)

	// It is ensured, that the passed command instance also
	// implements LimitedCommand.
	lcmd := cmd.(LimitedCommand)
	expireDuration := time.Duration(lcmd.LimiterBurst()) * lcmd.LimiterRestoration()

	limiter, ok := c.limiters.GetValue(key).(*Limiter)
	if ok {
		c.limiters.SetExpire(key, expireDuration)
		return limiter
	}

	limiter = c.pool.Get().(*Limiter).setParams(lcmd.LimiterBurst(), lcmd.LimiterRestoration())
	c.limiters.Set(key, limiter, expireDuration, func(val interface{}) {
		c.pool.Put(val)
	})

	return limiter
}
