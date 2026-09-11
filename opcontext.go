package mgo

import (
	"context"
	"time"
)

// OpTimeout bounds every MongoDB operation this package issues through the
// official driver.
//
// Before this existed each call passed context.Background(), so a query had no
// deadline and no cancellation. On a degraded replica set those calls do not
// give up: they accumulate, each holding a connection and pulling pages through
// the cache, which turns a slow secondary into a stuck one.
//
// Set to zero to restore the previous unbounded behaviour.
var OpTimeout = 30 * time.Second

// opContext returns a context bounded by OpTimeout.
//
// The caller must always call the returned cancel func. When an operation
// returns a cursor, the same context must stay alive for the whole iteration —
// cancel only after the cursor is drained, never between the query and All().
func opContext() (context.Context, context.CancelFunc) {
	if OpTimeout <= 0 {
		return context.Background(), func() {}
	}
	return context.WithTimeout(context.Background(), OpTimeout)
}
