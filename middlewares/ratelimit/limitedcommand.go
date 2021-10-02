package ratelimit

import "time"

// LimitedCommand specifies the structure of a
// rate limitable command.
type LimitedCommand interface {

	// GetLimiterBurst returns the maximum ammount
	// of tokens which can be available at a time.
	LimiterBurst() int

	// GetLimiterRestoration returns the duration
	// between new tokens are generated.
	LimiterRestoration() time.Duration

	// IsLimiterGlobal returns true if the limit
	// shall be handled globally across all guilds.
	// Otherwise, a limiter is created for each
	// guild the user executes the command on.
	IsLimiterGlobal() bool
}
