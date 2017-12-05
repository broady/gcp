package gcp

import (
	"context"
	"time"

	"golang.org/x/time/rate"

	"cloud.google.com/go/storage"
)

type ObjectUpdate struct {
	Err        error
	Generation int64
}

// WatchObject updates the channel whenever the content of a given object changes.
// If a limit is provided, it's used to rate limit the requests to GCS. If nil, WatchObject polls GCS every 5 seconds for updates.
func WatchObject(ctx context.Context, obj *storage.ObjectHandle, limit *rate.Limiter) chan ObjectUpdate {
	if limit == nil {
		limit = rate.NewLimiter(rate.Every(5*time.Second), 1)
	}
	update := make(chan ObjectUpdate)
	go func() {
		var currentMeta int64 = -1
		var oldError error
		for {
			select {
			case <-ctx.Done():
				update <- ObjectUpdate{Err: ctx.Err()}
				return
			default:
			}
			if err := limit.Wait(ctx); err != nil {
				update <- ObjectUpdate{Err: ctx.Err()}
				return
			}

			attrs, err := obj.Attrs(ctx)
			if err != nil {
				if oldError == nil || err.Error() != oldError.Error() {
					update <- ObjectUpdate{Err: err}
					oldError = err
				}
				continue
			}
			if currentMeta == attrs.Generation {
				continue
			}
			currentMeta = attrs.Generation
			update <- ObjectUpdate{Generation: attrs.Generation}
		}
	}()
	return update
}
